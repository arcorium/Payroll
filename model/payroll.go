package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payroll struct {
	Id                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	StaffSerialNumber uint16             `json:"no_urut_staff" bson:"no_urut_staff"`
	StaffName         string             `json:"nama_staff" bson:"nama_staff"`
	Institute         string             `json:"institusi" bson:"institusi"`
	Position          string             `json:"jabatan" bson:"jabatan"`
	Hours             uint8              `json:"jam_ajar" bson:"jam_ajar"`
	Salary            Salary             `json:"gaji" bson:"gaji"`
	Save              uint64             `json:"tabungan" bson:"tabungan"`
	SalaryCuts        SalaryCuts         `json:"potongan_gaji" bson:"potongan_gaji"`

	Total *Total `json:"total" bson:"total,omitempty"`

	Month uint8  `json:"months" bson:"months"`
	Years uint16 `json:"years" bson:"years"`
}

func (p *Payroll) SetTotal() {
	if p.Total == nil {
		p.Total = new(Total)
	}
	p.Total.Total = p.Salary.Total() - p.SalaryCuts.Total()
	p.Total.Receipt = p.Total.Total - p.Save
	p.Total.Slip = p.Total.Receipt + p.SalaryCuts.Nilam
}

type Salary struct {
	StaffSalary           uint64 `json:"gaji_pegawai" bson:"gaji_pegawai"`
	HonorarySalary        uint64 `json:"honor_ajar" bson:"honor_ajar"`
	HomeRoomTeacherSalary uint64 `json:"honor_walas" bson:"honor_walas"`
	TPMPS                 uint64 `json:"tpmps" bson:"tpmps"`
}

func (s *Salary) Total() uint64 {
	return s.StaffSalary + s.HonorarySalary + s.HomeRoomTeacherSalary + s.TPMPS
}

type SalaryCuts struct {
	Nilam     uint64 `json:"nilam" bson:"nilam"`
	BPJSTKSMP uint64 `json:"bpjs_tk_smp" bson:"bpjs_tk_smp"`
	BPJSSMP   uint64 `json:"bpjs_smp" bson:"bpjs_smp"`
	BPJSTKSMK uint64 `json:"bpjs_tk_smk" bson:"bpjs_tk_smk"`
	BPJSSMK   uint64 `json:"bpjs_smk" bson:"bpjs_smk"`
}

func (s *SalaryCuts) Total() uint64 {
	return s.Nilam + s.BPJSTKSMP + s.BPJSSMP + s.BPJSTKSMK + s.BPJSSMK
}

type HonorarySalary struct {
	SMP HonoraryDetails `json:"smp" bson:"smp"`
	SMK HonoraryDetails `json:"smk" bson:"smk"`
}

type HonoraryDetails struct {
	Hours uint8  `json:"jam" bson:"jam"`
	Total uint64 `json:"jumlah" bson:"jumlah"`
}

type Total struct {
	Total   uint64 `json:"total"`
	Receipt uint64 `json:"penerimaan"`
	Slip    uint64 `json:"slip"`
}

type PayrollRequest struct {
	Months uint8  `json:"months"`
	Years  uint16 `json:"years"`

	Data any `json:"data,omitempty"`
}
