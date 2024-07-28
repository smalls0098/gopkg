package gen

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"strings"
)

var DataTypeMap = map[string]func(gorm.ColumnType) (dataType string){
	"tinyint": func(columnType gorm.ColumnType) (dataType string) {
		ct, _ := columnType.ColumnType()
		if strings.HasSuffix(ct, "unsigned") {
			return "uint8"
		}
		return "bool"
	},
	"smallint": func(columnType gorm.ColumnType) (dataType string) {
		ct, _ := columnType.ColumnType()
		if strings.HasSuffix(ct, "unsigned") {
			return "uint32"
		}
		return "int32"
	},
	"mediumint": func(detailType gorm.ColumnType) (dataType string) {
		return "int64"
	},
	"int": func(columnType gorm.ColumnType) (dataType string) {
		ct, _ := columnType.ColumnType()
		if strings.HasSuffix(ct, "unsigned") {
			return "uint32"
		}
		return "int32"
	},
	"bigint": func(detailType gorm.ColumnType) (dataType string) {
		return "int64"
	},
}

var ModelOpt = []gen.ModelOpt{
	gen.FieldType("deleted_time", "soft_delete.DeletedAt"),
	gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
		tag.Set("autoUpdateTime")
		return tag
	}),
	gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
		tag.Set("autoCreateTime")
		return tag
	}),
	gen.FieldJSONTag("deleted_time", "-"),
}

var GenConfig = gen.Config{
	OutPath:      "./internal/data/dbquery",
	ModelPkgPath: "./internal/data/dbmodel",

	WithUnitTest: false,

	// generate model global configuration
	FieldNullable:     true,  // generate pointer when field is nullable
	FieldCoverable:    false, // generate pointer when field has default value
	FieldSignable:     true,  // detect integer field's unsigned type, adjust generated data type
	FieldWithIndexTag: true,  // generate with gorm index tag
	FieldWithTypeTag:  true,  // generate with gorm column type tag

	Mode: gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
}
