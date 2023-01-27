package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Staff struct {
	Id               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SerialNumber     uint16             `json:"no_urut" bson:"no_urut"`
	Name             string             `json:"nama" bson:"nama"`
	Status           string             `json:"status" bson:"status"`
	Position         InstitutePosition  `json:"jabatan" bson:"jabatan"`
	PersonalData     Personal           `json:"data_pribadi" bson:"data_pribadi"`
	Education        Education          `json:"pendidikan" bson:"pendidikan"`
	Savings          []Saving           `json:"tabungan,omitempty" bson:"tabungan"`
	TeachTimeDetails []TeachTimeDetail  `json:"details,omitempty" bson:"details"`
}

type InstitutePosition struct {
	Position        string          `json:"nama_jabatan" bson:"nama_jabatan"`
	Institute       string          `json:"instansi_induk" bson:"instansi_induk"`
	HomeRoomTeacher HomeRoomTeacher `json:"wali_kelas" bson:"wali_kelas"`
}

type HomeRoomTeacher struct {
	Class string `json:"kelas" bson:"kelas"`
	Type  string `json:"tipe" bson:"tipe"`
}

type Personal struct {
	BirthPlace string `json:"tempat_lahir" bson:"tempat_lahir"`
	//BirthDate  primitive.DateTime `json:"tanggal_lahir" bson:"tanggal_lahir"`
	BirthDate string `json:"tanggal_lahir" bson:"tanggal_lahir"`
	Address   string `json:"alamat" bson:"alamat"`
	Telephone string `json:"no_telepon" bson:"no_telepon"`
}

//type Address struct {
//	Road      string `json:"jalan" bson:"jalan"`
//	City      string `json:"kota" bson:"kota"`
//}

type Education struct {
	Education  string `json:"pendidikan_terakhir" bson:"pendidikan_terakhir"`
	SchoolName string `json:"nama_sekolah" bson:"nama_sekolah"`
	Major      string `json:"jurusan" bson:"jurusan"`
	Graduate   uint16 `json:"tahun_lulus" bson:"tahun_lulus"`
}

type TeachTimeDetail struct {
	UUID      string `json:"uuid,omitempty" bson:"uuid"`
	Institute string `json:"sekolah" bson:"sekolah"`
	Study     string `json:"pelajaran" bson:"pelajaran"`
	Hours     uint8  `json:"jumlah_jam" bson:"jumlah_jam"`

	Months uint8  `json:"month,omitempty" bson:"month"`
	Years  uint16 `json:"years,omitempty" bson:"years"`
}

type Saving struct {
	UUID  string `json:"uuid,omitempty" bson:"uuid"`
	Total uint64 `json:"total" bson:"total"`

	Months uint8  `json:"months" bson:"months"`
	Years  uint16 `json:"years" bson:"years"`
}
