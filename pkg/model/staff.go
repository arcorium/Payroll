package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Staff struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Jabatan     InstitutePosition  `json:"jabatan" bson:"jabatan"`
	DataPribadi Personal           `json:"data_pribadi" bson:"data_pribadi"`
	Pendidikan  Education          `json:"pendidikan" bson:"pendidikan"`
}

type InstitutePosition struct {
	Instansi  string          `json:"instansi" bson:"instansi"`
	Posisi    Position        `json:"posisi" bson:"posisi"`
	JamAjar   uint8           `json:"jam_ajar" bson:"jam_ajar"`
	WaliKelas HomeRoomTeacher `json:"wali_kelas" bson:"wali_kelas"`
}

type HomeRoomTeacher struct {
	Kelas uint8  `json:"kelas" bson:"kelas"`
	Abjad string `json:"abjad" bson:"abjad"`
}

type Personal struct {
	TempatLahir  string             `json:"tempat_lahir" bson:"tempat_lahir"`
	TanggalLahir primitive.DateTime `json:"tanggal_lahir" bson:"tanggal_lahir"`
	Address      Address            `json:"address" bson:"address"`
}

type Address struct {
	Jalan   string `json:"jalan" bson:"jalan"`
	Kota    string `json:"kota" bson:"kota"`
	Telepon string `json:"telepon" bson:"telepon"`
}

type Education struct {
	Pendidikan  string `json:"pendidikan" bson:"pendidikan"`
	NamaSekolah string `json:"nama_sekolah" bson:"nama_sekolah"`
	Jurusan     string `json:"jurusan" bson:"jurusan"`
	Lulusan     uint16 `json:"lulusan" bson:"lulusan"`
}
