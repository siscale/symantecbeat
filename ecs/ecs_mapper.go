package ecs

import (
	"encoding/csv"
	"github.com/elastic/beats/libbeat/logp"
	"io"
	"os"
	"strings"
)

type Mapper struct {
	fields map[string]string
	logger *logp.Logger
}

func NewMapper(file string) (*Mapper, error) {
	fields := make(map[string]string, 0)
	mapper := Mapper{
		logger: logp.NewLogger("ecs_mapper"),
		fields: fields,
	}

	err := mapper.readCsv(file)
	if err != nil {
		return nil, err
	}
	return &mapper, nil
}

func (m *Mapper) EcsField(oldFieldName string) string {
	if ecsField, ok := m.fields[oldFieldName]; ok {
		return ecsField
	}
	return oldFieldName
}

func (m *Mapper) Size() int {
	return len(m.fields)
}

func (m *Mapper) readCsv(file string) error {
	csvIn, err := os.Open(file)
	if err != nil {
		m.logger.Error("Error reading csv file for ecs_transformation err=%v", err)
		return err
	}
	r := csv.NewReader(csvIn)
	// handle header
	_, err = r.Read()
	if err != nil {
		m.logger.Error("Error reading csv header err=%v", err)
	}

	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			m.logger.Error("Error reading line from csv file err=%v", err)
		}
		oldField := rec[0]
		ecsField := rec[1]
		if len(strings.TrimSpace(ecsField)) != 0 {
			m.fields[oldField] = ecsField
		}
	}
	return nil
}
