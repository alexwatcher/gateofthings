package models

type Profile struct {
	Id            string
	IsProvisioned bool
	Name          string
	Avatar        []byte
}
