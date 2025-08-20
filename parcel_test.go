package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource — это источник псевдослучайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix-формате в виде числа
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// get
	stored, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.Equal(t, parcel, stored)

	// delete
	err = store.Delete(parcel.Number)

	stored, err = store.Get(parcel.Number)
	require.Equal(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)

	require.NoError(t, err)

	// check
	stored, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.Equal(t, newAddress, stored.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// set status
	err = store.SetStatus(parcel.Number, ParcelStatusSent)

	require.NoError(t, err)

	// check
	stored, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, stored.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам одного клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])

		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)

	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]

		require.True(t, ok)
		require.Equal(t, expectedParcel, parcel)
	}
}
