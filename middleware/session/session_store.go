package session

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaswangc/gowork/data/redis_cli"
	"github.com/asaswangc/gowork/utils"
	"github.com/asaswangc/gowork/utils/cfg"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
	"time"
)

const (
	sessionExpire  = 60 * 120 // Amount of time for cookies/redis keys to expire
	redisTTLAge    = 60 * 20  // 20 minutes
	redisMaxLength = 4096
)

// RedisClusterStore 自定义store结构体，用于实现gorilla store的相关方法
// 仿写自：https://github.com/boj/redistore/tree/v1.2
type RedisClusterStore struct {
	Pool          *redis.ClusterClient
	Codecs        []securecookie.Codec
	Options       *sessions.Options // default configuration
	DefaultMaxAge int               // default Redis TTL for a MaxAge == 0 session
	maxLength     int
	keyPrefix     string
}

func (s *RedisClusterStore) SetMaxLength(l int) {
	if l >= 0 {
		s.maxLength = l
	}
}

func (s *RedisClusterStore) SetKeyPrefix(p string) {
	s.keyPrefix = p
}

func (s *RedisClusterStore) SetMaxAge(v int) (err error) {
	var c *securecookie.SecureCookie
	var ok bool
	s.Options.MaxAge = v
	for i := range s.Codecs {
		if c, ok = s.Codecs[i].(*securecookie.SecureCookie); ok {
			c.MaxAge(v)
		} else {
			return fmt.Errorf("Can't change MaxAge on codec %v\n", s.Codecs[i])
		}
	}
	return err
}

func (s *RedisClusterStore) Close() error {
	return s.Pool.Close()
}

func (s *RedisClusterStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *RedisClusterStore) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.Options
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		id, err := utils.DecryptWithRSA(c.Value, cfg.T.KeyPath.Private)
		session.ID = id
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

func (s *RedisClusterStore) Save(_ *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		// 数据加密
		encryptedData, err := utils.EncryptWithRSA(session.ID, cfg.T.KeyPath.Public)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encryptedData, session.Options))
	}
	return nil
}

func (s *RedisClusterStore) load(session *sessions.Session) (bool, error) {
	res, err := s.Pool.Get(context.TODO(), s.keyPrefix+session.ID).Bytes()
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, nil // no data was associated with this key
	}
	return true, s.deserialize(res, session)
}

func (s *RedisClusterStore) delete(session *sessions.Session) error {
	return s.Pool.Del(context.TODO(), s.keyPrefix+session.ID).Err()
}

func (s *RedisClusterStore) save(session *sessions.Session) error {
	b, err := s.serialize(session)
	if err != nil {
		return err
	}
	if s.maxLength != 0 && len(b) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	age := session.Options.MaxAge
	if age == 0 {
		age = s.DefaultMaxAge
	}

	return s.Pool.SetEX(context.TODO(), s.keyPrefix+session.ID, b, time.Duration(age)*time.Second).Err()
}

func (s *RedisClusterStore) deserialize(d []byte, ss *sessions.Session) (err error) {
	m := make(map[string]interface{})
	if err = json.Unmarshal(d, &m); err != nil {
		return fmt.Errorf("反序列化失败,Error:%s", err)
	}
	for k, v := range m {
		ss.Values[k] = v
	}
	return err
}

func (s *RedisClusterStore) serialize(ss *sessions.Session) (r []byte, err error) {
	m := make(map[string]interface{}, len(ss.Values))
	for k, v := range ss.Values {
		ks, ok := k.(string)
		if !ok {
			err := fmt.Errorf("non-string key value, cannot serialize session to JSON: %v", k)
			return r, fmt.Errorf("反序列化失败,Error:%s", err)
		}
		m[ks] = v
	}
	return json.Marshal(m)
}

func NewRedisClusterStore(keyPrefix string, redisAge, maxLength int, secure, httpOnly bool, sameSiteMode http.SameSite, _ ...[]byte) (*RedisClusterStore, error) {
	if redisAge == 0 {
		redisAge = redisTTLAge
	}
	if maxLength == 0 {
		maxLength = redisMaxLength
	}
	if sameSiteMode == 0 {
		sameSiteMode = http.SameSiteDefaultMode
	}

	rs := &RedisClusterStore{
		Pool:   redis_cli.RedisDB,
		Codecs: securecookie.CodecsFromPairs(),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   sessionExpire,
			Secure:   secure,
			HttpOnly: httpOnly,
			SameSite: sameSiteMode,
		},
		DefaultMaxAge: redisAge,
		maxLength:     maxLength,
		keyPrefix:     keyPrefix,
	}
	return rs, nil
}
