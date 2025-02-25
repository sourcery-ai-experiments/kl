package client

import (
	"fmt"
	"strings"

	fn "github.com/kloudlite/kl/pkg/functions"
)

type ResEnvType struct {
	Name   string `json:"name" yaml:"name"`
	Key    string `json:"key"`
	RefKey string `json:"refKey" yaml:"refKey"`
}

type EnvType struct {
	Key       string  `json:"key" yaml:"key"`
	Value     *string `json:"value,omitempty" yaml:"value,omitempty"`
	ConfigRef *string `json:"configRef,omitempty" yaml:"configRef,omitempty"`
	SecretRef *string `json:"secretRef,omitempty" yaml:"secretRef,omitempty"`
	MresRef   *string `json:"mresRef,omitempty" yaml:"mresRef,omitempty"`
}

type ResType struct {
	Name string       `json:"name"`
	Env  []ResEnvType `json:"env"`
}

type EnvVars []EnvType

type NormalEnv struct {
	Key   string
	Value string
}

func (e *EnvVars) GetEnvs() []NormalEnv {
	resp := make([]NormalEnv, 0)
	if e == nil {
		return resp
	}

	for _, r := range *e {
		if r.Value != nil {
			resp = append(resp, NormalEnv{
				Key:   r.Key,
				Value: *r.Value,
			})
		}
	}

	return resp
}

type resType string

const (
	Res_config resType = "config"
	Res_secret resType = "secret"
	Res_mres   resType = "mres"
)

func (e *EnvVars) getReses(res resType) []ResType {

	resp := make([]ResType, 0)
	if e == nil {
		return resp
	}

	hist := map[string]int{}

	for _, r := range *e {
		var ref *string

		switch res {
		case Res_config:
			ref = r.ConfigRef
		case Res_secret:
			ref = r.SecretRef
		case Res_mres:
			ref = r.MresRef
		default:
			continue
		}

		if ref == nil {
			continue
		}

		s := strings.Split(*ref, "/")
		if len(s) != 2 {
			continue
		}

		mName, mKey := s[0], s[1]

		j, ok := hist[mName]
		if !ok {
			hist[mName] = len(resp)
			resp = append(resp, ResType{
				Name: mName,
				Env: []ResEnvType{
					{
						Key:    r.Key,
						RefKey: mKey,
					},
				},
			})
			continue
		}

		resp[j].Env = append(resp[j].Env, ResEnvType{
			Key:    r.Key,
			RefKey: mKey,
		})
	}

	return resp
}

func (e *EnvVars) GetMreses() []ResType {
	return e.getReses(Res_mres)
}

func (e *EnvVars) GetConfigs() []ResType {
	return e.getReses(Res_config)
}

func (e *EnvVars) GetSecrets() []ResType {
	return e.getReses(Res_secret)
}

func (e *EnvVars) AddResTypes(rt []ResType, rtype resType) {

	if e == nil {
		e = &EnvVars{}
	}

	keys := map[string]bool{}

	getEnvKey := func(r EnvType) string {
		return fmt.Sprint(r.Key, func() string {
			if r.SecretRef != nil {
				return *r.SecretRef
			}
			if r.MresRef != nil {
				return *r.MresRef
			}
			if r.SecretRef != nil {
				return *r.SecretRef
			}
			if r.Value != nil {
				return *r.Value
			}

			return ""
		}())
	}

	getRtKey := func(name, key, refKey string) string {
		return fmt.Sprint(key, name, "/", refKey)
	}

	for _, r := range *e {
		ek := getEnvKey(r)

		if !keys[ek] {
			keys[ek] = true
		}
	}

	appendEnv := func(key, name, refKey string) {
		*e = append(*e, EnvType{
			Key:   key,
			Value: nil,
			ConfigRef: func() *string {
				if rtype != Res_config {
					return nil
				}

				return fn.Ptr(fmt.Sprint(name, "/", refKey))
			}(),
			SecretRef: func() *string {
				if rtype != Res_secret {
					return nil
				}

				return fn.Ptr(fmt.Sprint(name, "/", refKey))
			}(),
			MresRef: func() *string {
				if rtype != Res_mres {
					return nil
				}

				return fn.Ptr(fmt.Sprint(name, "/", refKey))
			}(),
		})
	}

	for _, r := range rt {
		for _, ret := range r.Env {
			ek := getRtKey(r.Name, ret.Key, ret.RefKey)
			if !keys[ek] {
				keys[ek] = true
				appendEnv(ret.Key, r.Name, ret.RefKey)
			}
		}
	}
}
