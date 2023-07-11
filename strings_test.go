package sqlb

import (
	"testing"
)

func TestCompact(t *testing.T) {
	src := `
		SELECT
			*
		FROM
			table
	`
	got := Compact(src)
	want := "SELECT * FROM table"
	if got != want {
		t.Errorf("%s = %q, want %q", "test compact #1", got, want)
	}

	src = ""
	got = Compact(src)
	want = ""
	if got != want {
		t.Errorf("%s = %q, want %q", "test compact #2", got, want)
	}

	src = "SELECT"
	got = Compact(src)
	want = "SELECT"
	if got != want {
		t.Errorf("%s = %q, want %q", "test compact #3", got, want)
	}
}

func BenchmarkCompact(b *testing.B) {
	const query = `INSERT INTO #rac.ReserveT (
		DelFlg,
		UpdateDateTime,
		InsertDateTime,
		ProcessID,
		ProcessPGM,
		ProcessNo,
		ShopID,
		CustID,
		UseCustID,
		ClassID,
		CarID,
		Status,
		ReserveDiv,
		ChargeType,
		ReserveType,
		ReservePurpose,
		RegistDate,
		RegistPersonID,
		StartDateTime,
		ReturnDateTime,
		DiffTime,
		DiffDay,
		DriveType,
		Smoking,
		OptionNavi,
		OptionSeat,
		OptionSeatCnt,
		OptionPet,
		OptionStud,
		OptionETC,
		OptionPick,
		OptionDeli,
		Insurance,
		Biko,
		ChargeTotal,
		Receipt,
		CompanyID,
		GroupID,
		ReturnDateTimeOFF,
		PersonID,
		FreeCustName,
		FreeCustDeptName,
		FreeCustPersonName,
		FreeUseCustName,
		FreeUseCustDeptName,
		FreeUseCustPersonName,
		FreeDeliveryName,
		FreeDeliveryTEL,
		FreeDeliveryAddress,
		OptionOther,
		StartDelayDateTime,
		ReturnDelayDateTime
	)
	SELECT
		0 AS DelFlg,
		@now AS UpdateDateTime,
		@now AS InsertDateTime,
		0 AS ProcessID,
		'reserve.v3.0' AS ProcessPGM,
		1 AS ProcessNo,
		ShopM.ShopID AS ShopID,
		@custID AS CustID,
		0 AS UseCustID,
		@classID AS ClassID,
		@carID AS CarID,
		1 AS Status,
		1 AS ReserveDiv,
		ShopM.ChargeType AS ChargeType,
		3 AS ReserveType,
		IF(FcGroupM.SystemMode = 1, 12, 0) AS ReservePurpose,
		now() AS RegistDate,
		0 AS RegistPersonID,
		@beginTime AS StartDateTime,
		@endTime AS ReturnDateTime,
		HOUR(TIMEDIFF(@endTime, @beginTime)) + IF(MINUTE(TIMEDIFF(@endTime, @beginTime)) > 0, 1, 0) AS DiffTime,
		DATEDIFF(@endTime, @beginTime) + 1 AS DiffDay,
		CarM.DriveType AS DriveType,
		CarM.Smoking AS Smoking,
		IF(@navi > 0, @navi, CarM.CarOptionNavi) AS OptionNavi,
		IF(@seat > 0, 1, 0) AS OptionSeat,
		@seat AS OptionSeatCnt,
		@pet AS OptionPet,
		@studless AS OptionStud,
		CarM.CarOptionETC AS OptionETC,
		@pick AS OptionPick,
		@deli AS OptionDeli,
		@insurance AS Insurance,
		@note AS Biko,
		@charge_total AS ChargeTotal,
		@charge_total AS Receipt,
		ShopM.CompanyID AS CompanyID,
		ShopM.GroupID AS GroupID,
		0 AS ReturnDateTimeOFF,
		0 AS PersonID,
		'' AS FreeCustName,
		'' AS FreeCustDeptName,
		'' AS FreeCustPersonName,
		@drivers AS FreeUseCustName,
		'' AS FreeUseCustDeptName,
		'' AS FreeUseCustPersonName,
		'' AS FreeDeliveryName,
		'' AS FreeDeliveryTEL,
		'' AS FreeDeliveryAddress,
		@option_other AS OptionOther,
		@beginTime AS StartDelayDateTime,
		@endTime AS ReturnDelayDateTime
	FROM
		#rac.ShopM INNER JOIN
		#rac.CompanyM ON (
			CompanyM.CompanyID = ShopM.CompanyID
		) INNER JOIN
		#rac.FcGroupM ON (
			FcGroupM.FcGroupID = CompanyM.FcGroupID
		),
		#rac.CarM
	WHERE
		ShopM.ShopID = @shopID
	AND CarM.CarID = @carID
	`

	for i := 0; i < b.N; i++ {
		Compact(query)
	}
}

// BenchmarkCompact-4   	  328986	      3361 ns/op	    2688 B/op	       1 allocs/op
