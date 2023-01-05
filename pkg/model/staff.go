package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Staff struct {
	Id               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Position         InstitutePosition  `json:"jabatan" bson:"jabatan"`
	PersonalData     Personal           `json:"data_pribadi" bson:"data_pribadi"`
	Education        Education          `json:"pendidikan" bson:"pendidikan"`
	TeachTimeDetails []TeachTimeDetail  `json:"details" bson:"details"`
}

type InstitutePosition struct {
	Institute       string          `json:"instansi" bson:"instansi"`
	Position        Position        `json:"posisi" bson:"posisi"`
	TeachTime       uint8           `json:"jam_ajar" bson:"jam_ajar"`
	HomeRoomTeacher HomeRoomTeacher `json:"wali_kelas" bson:"wali_kelas"`
}

type HomeRoomTeacher struct {
	Class uint8  `json:"kelas" bson:"kelas"`
	Type  string `json:"tipe" bson:"tipe"`
}

type Personal struct {
	BirthPlace string             `json:"tempat_lahir" bson:"tempat_lahir"`
	BirthDate  primitive.DateTime `json:"tanggal_lahir" bson:"tanggal_lahir"`
	Address    Address            `json:"address" bson:"address"`
}

type Address struct {
	Road      string `json:"jalan" bson:"jalan"`
	City      string `json:"kota" bson:"kota"`
	Telephone string `json:"telepon" bson:"telepon"`
}

type Education struct {
	Education  string `json:"pendidikan" bson:"pendidikan"`
	SchoolName string `json:"nama_sekolah" bson:"nama_sekolah"`
	Major      string `json:"jurusan" bson:"jurusan"`
	Graduate   uint16 `json:"lulusan" bson:"lulusan"`
}

type TeachTimeDetail struct {
	Institute string `json:"sekolah" bson:"sekolah"`
	Study     string `json:"pelajaran" bson:"pelajaran"`
	Hours     uint8  `json:"jumlah_jam" bson:"jumlah_jam"`
}
