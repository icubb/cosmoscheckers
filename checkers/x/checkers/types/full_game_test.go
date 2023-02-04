package types_test

const (
	alice = testutil.Alice
	bob   = testutil.Bob
)

func GetStoredGame1() *types.StoredGame {
	return types.StoredGame{
		Black: alice,
		Red:   bob,
		Index: "1",
		Board: rules.New().String(),
		Turn:  "b",
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	aliceAddress, err1 := sdk.AccAddressFromBech32(alice)
	black, err2 := GetStoredGame1().GetBlackAddress()
	require.Equal(t, aliceAddress, black)
	require.Nil(t, err2)
	require.Nil(t, err1)
}

// Requires that alice address is equal to the black one
// Requires that the errors are nil.

func TestGetAddressWrongBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4" //bad last digit
	black, err := storedGame.GetBlackAddress()
	require.Nil(t, black)
	require.EqualError(t,
		err,
		"black address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

// so the black is required to be nil since the getblackaddress method should fail on the stored game

func TestParseDeadlineCorrect(t *testing.T) {
	deadline, err := GetStoredGame1().GetDeadlineAsTime()
	require.Nil(t, err)
	require.Equal(t, time.Time(time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.UTC)), deadline)
}

func TestParseDeadlineMissingMonth(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Deadline = "2006-02 15:04:05.999999999 +0000 UTC"
	_, err := storedGame.GetDeadlineAsTime()
	require.EqualError(t,
		err,
		"deadline cannot be parsed: 2006-02 15:04:05.999999999 +0000 UTC: parsing time \"2006-02 15:04:05.999999999 +0000 UTC\" as \"2006-01-02 15:04:05.999999999 +0000 UTC\": cannot parse \" 15:04:05.999999999 +0000 UTC\" as \"-\"")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestGameValidateOk(t *testing.T) {
	storedGame := GetStoredGame1()
	require.NoError(t, storedGame.Validate())
}

func TestGetPlayerAddressBlackCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	black, found, err := storedGame.GetPlayerAddress("b")
	require.Equal(t, alice, black.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressBlackIncorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "notanaddress"
	black, found, err := storedGame.GetPlayerAddress("b")
	require.Nil(t, black)
	require.False(t, found)
	require.EqualError(t, err, "black address is invalid: not an address: decoding bech32 failed: invalid separator index -1")
}

func TestGetPlayerAddressRedCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	red, found, err := storedGame.GetPlayerAddress("r")
	require.Equal(t, bob, red.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressBlackIncorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "notanaddress"
	black, found, err := storedGame.GetPlayerAddress("r")
	require.Nil(t, red)
	require.False(t, found)
	require.EqualError(t, err, "red address is invalid: not an address: decoding bech32 failed: invalid separator index -1")
}

func TestGetPlayerAddressAnyNotFound(t *testing.T) {
	storedGame := GetStoredGame1()
	white, found, err := storedGame.GetPlayerAddress("*")
	require.Nil(t, white)
	require.False(t, found)
	require.Nil(t, err)
}

func TestGetWinnerBlackCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Winner = "b"
	winner, found, err := storedGame.GetWinnerAddress()
	require.Equal(t, alice, winner.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetWinnerRedCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Winner = "r"
	winner, found, err := storedGame.GetWinnerAddress()
	require.Equal(t, bob, winner.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetWinnerNotYetCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	winner, found, err := storedGame.GetWinnerAddress()
	require.Nil(t, winner)
	require.False(t, found)
	require.Nil(t, err)
}
