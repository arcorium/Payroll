package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payroll struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	StaffId    primitive.ObjectID `json:"staff_id" bson:"staff_id"`
	StaffName  string             `json:"nama_staff,omitempty" bson:"nama_staff,omitempty"`
	Salary     Salary             `json:"gaji" bson:"gaji"`
	SalaryCuts SalaryCuts         `json:"potongan_gaji" bson:"potongan_gaji"`
	Month      uint8              `json:"bulan" bson:"bulan"`
	Years      uint16             `json:"tahun" bson:"tahun"`
}

type Salary struct {
	StaffSalary           uint64         `json:"gaji_pegawai" bson:"gaji_pegawai"`
	HonorarySalary        HonorarySalary `json:"honor_ajar" bson:"honor_ajar"`
	HomeRoomTeacherSalary uint64         `json:"honor_walas" bson:"honor_walas"`
	TPMPS                 uint64         `json:"tpmps" bson:"tpmps"`
}

type SalaryCuts struct {
	Save      uint64 `json:"tabungan" bson:"tabungan"`
	Nilam     uint64 `json:"nilam" bson:"nilam"`
	BPJSTKSMP uint64 `json:"bpjs_tk_smp" bson:"bpjs_tk_smp"`
	BPJSSMP   uint64 `json:"bpjs_smp" bson:"bpjs_smp"`
	BPJSTKSMK uint64 `json:"bpjs_tk_smk" bson:"bpjs_tk_smk"`
	BPJSSMK   uint64 `json:"bpjs_smk" bson:"bpjs_smk"`
}

type HonorarySalary struct {
	SMP HonoraryDetails `json:"smp" bson:"smp"`
	SMK HonoraryDetails `json:"smk" bson:"smk"`
}

type HonoraryDetails struct {
	Hours uint8  `json:"jam" bson:"jam"`
	Total uint64 `json:"jumlah" bson:"jumlah"`
}
