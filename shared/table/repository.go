package table

import (
	"errors"
	"path"

	"github.com/rs/zerolog/log"
)

var ErrTableRepositoryIsNil = errors.New("table repository is nil")

type Repository struct {
	// 아이템 테이블
	Item ItemTable

	// 캐릭터 생성시 장비 확률 테이블
	CreateCharacter CreateCharacterTable

	// 캐릭터 생성시 한손 무기 확률 테이블
	CreateCharacter1HWeapon CreateCharacter1HWeaponTable

	// 캐릭터 생성시 두손 무기 확률 테이블
	CreateCharacter2HWeapon CreateCharacter2HWeaponTable

	// 캐릭터 생성시 모자 확률 테이블
	CreateCharacterCap CreateCharacterCapTable

	// 캐릭터 생성시 망토 확률 테이블
	CreateCharacterCape CreateCharacterCapeTable

	// 캐릭터 생성시 상의 확률 테이블
	CreateCharacterCoat CreateCharacterCoatTable

	// 캐릭터 생성시 귀 형태 확률 테이블
	CreateCharacterEar CreateCharacterEarTable

	// 캐릭터 생성시 귀 장식 확률 테이블
	CreateCharacterEarAcc CreateCharacterEarAccTable

	// 캐릭터 생성시 눈 장식 확률 테이블
	CreateCharacterEysAcc CreateCharacterEysAccTable

	// 캐릭터 생성시 얼굴 확률 테이블
	CreateCharacterFace CreateCharacterFaceTable

	// 캐릭터 생성시 얼굴 장식 확률 테이블
	CreateCharacterFaceAcc CreateCharacterFaceAccTable

	// 캐릭터 생성시 장갑 확률 테이블
	CreateCharacterGlove CreateCharacterGloveTable

	// 캐릭터 생성시 헤어 확률 테이블
	CreateCharacterHair CreateCharacterHairTable

	// 캐릭터 생성시 한벌옷 확률 테이블
	CreateCharacterLongCoat CreateCharacterLongCoatTable

	// 캐릭터 생성시 하의 확률 테이블
	CreateCharacterPants CreateCharacterPantsTable

	// 캐릭터 생성시 신발 확률 테이블
	CreateCharacterShoes CreateCharacterShoesTable

	// 캐릭터 생성시 피부 확률 테이블
	CreateCharacterSkin CreateCharacterSkinTable

	// 캐릭터 생성시 서브 무기 확률 테이블
	CreateCharacterSubWeapon CreateCharacterSubWeaponTable
}

func (r *Repository) Load(csvDirPath string) error {
	var errs error

	r.Item = *NewItemTable()
	if err := r.Item.Load(path.Join(csvDirPath, "Item.csv")); err != nil {
		log.Err(err).Msg("failed to load Item table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacter = *NewCreateCharacterTable()
	if err := r.CreateCharacter.Load(path.Join(csvDirPath, "CreateCharacter.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacter table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacter1HWeapon = *NewCreateCharacter1HWeaponTable()
	if err := r.CreateCharacter1HWeapon.Load(path.Join(csvDirPath, "CreateCharacter1HWeapon.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacter1HWeapon table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacter2HWeapon = *NewCreateCharacter2HWeaponTable()
	if err := r.CreateCharacter2HWeapon.Load(path.Join(csvDirPath, "CreateCharacter2HWeapon.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacter2HWeapon table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterCap = *NewCreateCharacterCapTable()
	if err := r.CreateCharacterCap.Load(path.Join(csvDirPath, "CreateCharacterCap.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterCap table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterCape = *NewCreateCharacterCapeTable()
	if err := r.CreateCharacterCape.Load(path.Join(csvDirPath, "CreateCharacterCape.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterCape table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterCoat = *NewCreateCharacterCoatTable()
	if err := r.CreateCharacterCoat.Load(path.Join(csvDirPath, "CreateCharacterCoat.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterCoat table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterEar = *NewCreateCharacterEarTable()
	if err := r.CreateCharacterEar.Load(path.Join(csvDirPath, "CreateCharacterEar.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterEar table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterEarAcc = *NewCreateCharacterEarAccTable()
	if err := r.CreateCharacterEarAcc.Load(path.Join(csvDirPath, "CreateCharacterEarAcc.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterEarAcc table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterEysAcc = *NewCreateCharacterEysAccTable()
	if err := r.CreateCharacterEysAcc.Load(path.Join(csvDirPath, "CreateCharacterEysAcc.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterEysAcc table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterFace = *NewCreateCharacterFaceTable()
	if err := r.CreateCharacterFace.Load(path.Join(csvDirPath, "CreateCharacterFace.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterFace table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterFaceAcc = *NewCreateCharacterFaceAccTable()
	if err := r.CreateCharacterFaceAcc.Load(path.Join(csvDirPath, "CreateCharacterFaceAcc.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterFaceAcc table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterGlove = *NewCreateCharacterGloveTable()
	if err := r.CreateCharacterGlove.Load(path.Join(csvDirPath, "CreateCharacterGlove.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterGlove table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterHair = *NewCreateCharacterHairTable()
	if err := r.CreateCharacterHair.Load(path.Join(csvDirPath, "CreateCharacterHair.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterHair table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterLongCoat = *NewCreateCharacterLongCoatTable()
	if err := r.CreateCharacterLongCoat.Load(path.Join(csvDirPath, "CreateCharacterLongCoat.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterLongCoat table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterPants = *NewCreateCharacterPantsTable()
	if err := r.CreateCharacterPants.Load(path.Join(csvDirPath, "CreateCharacterPants.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterPants table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterShoes = *NewCreateCharacterShoesTable()
	if err := r.CreateCharacterShoes.Load(path.Join(csvDirPath, "CreateCharacterShoes.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterShoes table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterSkin = *NewCreateCharacterSkinTable()
	if err := r.CreateCharacterSkin.Load(path.Join(csvDirPath, "CreateCharacterSkin.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterSkin table")
		errs = errors.Join(errs, err)
	}

	r.CreateCharacterSubWeapon = *NewCreateCharacterSubWeaponTable()
	if err := r.CreateCharacterSubWeapon.Load(path.Join(csvDirPath, "CreateCharacterSubWeapon.csv")); err != nil {
		log.Err(err).Msg("failed to load CreateCharacterSubWeapon table")
		errs = errors.Join(errs, err)
	}

	return errs
}
