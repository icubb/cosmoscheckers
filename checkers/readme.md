# checkers
**checkers** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Ignite CLI has scaffolded a Vue.js-based web app in the `vue` directory. Run the following commands to install dependencies and start the app:

```
cd vue
npm install
npm run serve
```

The frontend app is built using the `@starport/vue` and `@starport/vuex` packages. For details, see the [monorepo for Ignite front-end development](https://github.com/ignite-hq/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/alice/checkers@latest! | sudo bash
```
`alice/checkers` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)


In this section you will handle:
    - The stored game object 
    - Protobuf objects 
    - Query.proto
    - Protobuf service interfaces
    - Your first unit test
    - Interactions via the command-line


In the Ignite CLI introduction section you learned how to start a completely new blockchain. Now it is time to dive deeper and explore how you can create a blockchain to play a decentralized game of checkers.


**Some Initial Thoughts** 

- As you are face to face with the proverbial blank page: where do you start? 

- A good place to start is thinking about the objects you keep in storage. **A game** obviously... but what does any game have to keep in storage?

- Questions to ask that could influence your design include, and are not limited to:
    - What is the lifecycle of a game?
    - How are the participants selected to be in the game?
    - What fields make it possible to differentiate between different games?
    - How do you ensure saftey against malice, sabotage, or even simple errors?
    - What limitations does your design **intentionally** impose on participants?
    - What limitations does your design **unintentionally** impose on participants?

- After thinking about what goes into each individual game, you should consider the demands of the wider system. In general terms, before you think about the commands that acheive what you seek, ask:
    - How do you lay games in storage?
    - How do you save and retrieve games? 

- The goal here is not to finalize every conceivable game feature immediately. For instance, handling wagers or leaderboards can be left for another time.
- But there should be a basic game design good enough to accommodate future improvements. 


**Code Needs** 

- **Do not** dive headlong into coding the rules of Checkers in go - examples will already exist which you can put to use. Your job is to make a blockchain that happens to enable the game of checkers. 

- With that in mind:
    - What Ignite CLI commands will get you a long way when it comes to implementation?
    - How do you adjust what Ignite CLI created for you?
    - How would you unit-test your modest additions?
    - How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get? 

- Run the commands, make the adjustments, run some tests regarding game storage. Do not go into deeper issues like messages and transactions yet. 

**Defining the rule set** 

- https://tutorials.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html

- A good start to developing a checkers blockchain is to define the rule set of the game. There are many versions of the rules. Choose a very simple set of basic rules to avoid getting lost in the rules of checkers or the proper implementation of the board state. 

- Use a ready-made implementation (opens new window) with the additional rule that the board is 8x8, is played on black cells, and black plays first. This code will not need adjustments. Copy this rules file into a rules folder inside your module. Change its package from checkers to rules. You can do this by command-line:

- Do not focus on the GUI, this procedure lays the foundation for an interface.

- Now it's time to create the first object.


**The stored game object** 

- Begin with the minimum game information needed to be stored:
    - **Black Player**: A string, the serialized address.
    - **Red Player**: A string, the serialized address.
    - **Board proper**: A string, the board as it is serialized by the *rules* file.
    - **Player to play next**: A string, specifiying whose *turn* it is.

- When you save strings, it makes it easier to understand what comes straight out of storage, but at the expense of storage space. 
- As an advanced consideration, you could store the same information in binary. 

**How to store** 

- After you know **what** to store, you have to decide **how** to store a game. This is important if you want your blockchain application to accommodate multiple simultaneous games.

- The game is identified by a unique ID.

- How should you generate the ID? If you let players choose it themselves, this could lead to transactions failing because of an ID clash. You cannot rely on a large random number like a universally unique identifier (UUID), because transactions have to be verifiable in the future. Verifiable means that nodes verifying the block need to arrive at the same conclusion. However, the new UUID() command is not deterministic. In this context, it is better to have a counter incrementing on each new game. This is possible because the code execution happens in a single thread.

- The counter must be kept in storage between transactions. Instead of a single counter in storage, you can keep the counter in a unique object at a singluar storage location, and easily add relevant elements to the objects as needed in the future. Name the counter as `nextId` and its container as `SystemInfo`.

- As for the game type, you can name it as `StoredGame`.

- You can rely on Ignite CLI's assistance for both the counter and the game. 


$ docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold single systemInfo nextId:uint \
    --module checkers \
    --no-message

`docker exec -it checkers ignite scaffold single systemInfo nextId:uint --module checkers --no-message`

- In this command
    - The `nextId` is explicitly made to be a `uint`. If you left it to Ignite's default, it would be a `string`.

    - You must add `--no-message`. If you omit it, Ignite CLI creates an `sdk.Msg` and an associated service whose purpose is to overwrite your `SystemInfo` object. However, your `SystemInfo.NextId` must be controlled/incremented by the application and not by the player sending a value of their own choosing. Ignite CLI still creates convenient getters. 

    So having no messgage means the people cannot change the value but only the application can.

- For the game type, because you are storing games by ID, you need a map. Instruct Ignite CLI with `scaffold map` using the `StoredGame` name:

```
$ docker run --rm -it \
    -v $(pwd):/checkers \
    -w /checkers \
    checkers_i \
    ignite scaffold map storedGame board turn black red \
    --index index \
    --module checkers \
    --no-message
```

- In this command:
    - `board`, `turn`, `black` and `red` are by default strings, so there is no need to be explicit with for instance `board:string`.

    - `index` is the id field picked, and anyway is the default name when scaffolding a map. `id` cannot be chosen when scaffolding with ignite. 

    - `--no-message` prevents game objects from being created or overwritten with a simple `sdk.Msg`. The application instead creates and updates the objects when creating properly crafted messages like create game or play a move.

**Looking Around** 

- The command added new constants:

```
const (
    SystemInfoKey = "SystemValue-value-"
)
```

```
const (
    StoredGameKeyPrefix = "StoredGame/value/"
)
```

`In that file it continues`

```
cosnt (
    StoredGameKeyPrefix = "StoredGame/value"
)

// StoredGameKey returns the store key to retrieve a stored game from the index fields

func StoredGameKey (
    index string,
) []byte {
    var key []byte

    indexBytes := []byte(index)
    key = append(key, indexBytes...)
    key = append(key, []byte("/")...)

    return key
}
```



- These constants are used as prefixes for the keys that can access the storage location of objects.

- In the case of games, the store model lets you *narrow* the search. For instance:

https://tutorials.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html

`store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.StoredGameKeyPrefix))`

- This gets the store to access any game if you have its index. 

`b := store.Get(types.StoredGameKey(index,))`

Gets the store to access the game store.
stored game key essentially parses the index so it can be retrieved by the get function on store.

Need to understadn the new ctx.kvstore stuff better 

The Cosmos SDK comes with a large set of stores to persist the state of applications. By default, the main store of Cosmos SDK applications is a `multistore`, i.e. a store of stores. Developers can add any number of key-value stores to the multistore, depending on their application needs. The multistore exists to support the modularity of the Cosmos SDK, as it lets each module declare and manage their own subset of the state Key-value stores in the multistore can only be accessed with a specific capability `key`, which is typically held in the `keeper` of the module that declared the store. 



**Protobuf Objects** 


Ignite CLI creates the Protobuf objects in the `proto` directory before compiling them. The `SystemInfo` object looks like this:

```
message SystemInfo {
    uint64 nextId = 1;
}
```

- The `StoredGame` object looks like this:

```
message StoredGame {
    string index = 1;
    string board = 2;
    string turn = 3; 
    string black = 4; 
    string red = 5; 
}
```

- Both objects compile to:

```
type SystemInfo struct {
    NextId uint64 `protobuf:"varint,1,opt,name=nextId,proto3" json:"nextId,omitempty"`
}
```

- And to:

```
type StoredGame struct {
    Index string `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
    Board string `protobuf:"bytes,2,opt,name=board,proto3" json:"board,omitempty"`
    Turn  string `protobuf:"bytes,3,opt,name=turn,proto3" json:"turn,omitempty"`
    Black string `protobuf:"bytes,4,opt,name=black,proto3" json:"black,omitempty"`
    Red   string `protobuf:"bytes,5,opt,name=red,proto3" json:"red,omitempty"`
}
```

- These are not the only created Protobuf objects. The genesis state is also defined in Protobuf:

```
import "checkers/system_info.proto";
import "checkers/stored_game.proto";

message GenesisState {
    ...
    SystemInfo systemInfo = 2;
    repeated StoredGame storedGameList = 3 [(gogoproto.nullable) = false];
}
```

- This is compiled to:

```
type GenesisState struct {
    Params         Params       `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
    SystemInfo     *SystemInfo  `protobuf:"bytes,2,opt,name=systemInfo,proto3" json:"systemInfo,omitempty"`
    StoredGameList []StoredGame `protobuf:"bytes,3,rep,name=storedGameList,proto3" json:"storedGameList"`
}
```

- You can find query objects as part of the boilerplate objects created by Ignite CLI. Ignite CLI creates the objects according to a model, but this does not prevent you from making changes later if you decide these queries are not needed: 

```
message QueryGetSystemInfoRequest {}

message QueryGetSystemInfoResponse {
	SystemInfo SystemInfo = 1 [(gogoproto.nullable) = false];
}

```

- The query objects for `StoredGame` are more useful for your checkers game, and look like this:

```
message QueryGetStoredGameRequest {
    string index = 1;
}

message QueryGetStoredGameResponse {
    StoredGame StoredGame = 1 [(gogoproto.nullable) = false];
}

message QueryAllStoredGameRequest {
    cosmos.base.query.v1beta1.PageRequest pagination = 1;
}

message QueryAllStoredGameResponse {
    repeated StoredGame StoredGame = 1 [(gogoproto.nullable) = false];
    cosmos.base.query.v1beta1.PageResponse pagination = 2;
}

Then set this code


```

    const DefaultIndex uint64 = 1

    func DefaultGenesis() *GenesisState {
        return &GenesisState{
-          SystemInfo:     nil,
+          SystemInfo: SystemInfo{
+              NextId: uint64(DefaultIndex),
+          },
            StoredGameList: []StoredGame{},
            ...
        }
    }
```


- You can choose to start with no games or insert a number of games to start with. In either case, you must choose the first ID of the first future created game, which here is set at 1 by reusing the DefaultIndex value.

**Protobuf Service Interfaces** 

- In addition to created objects, ignite CLI also creates services that declare and define how to access newly created storage objects. Ignite CLI introduces empty service interfaces that can be filled as you add objects and messages when scaffolding a brand new module. 

- In this case, Ignite CLI added to service Query how to query for your objects:

```
service Query { 

    rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
        option (google.api.http).get = "alice/checkers/checkers/params";
    }

    rpc SystemInfo(QueryGetSystemInfoRequest) returns (QueryGetSystemInfoResponse) {
		option (google.api.http).get = "/alice/checkers/checkers/system_info";
	}

	rpc StoredGame(QueryGetStoredGameRequest) returns (QueryGetStoredGameResponse) {
		option (google.api.http).get = "/alice/checkers/checkers/stored_game/{index}";
	}

	rpc StoredGameAll(QueryAllStoredGameRequest) returns (QueryAllStoredGameResponse) {
		option (google.api.http).get = "/alice/checkers/checkers/stored_game";
	}
}
```

Ignite CLI separates concerns into different files in the compilation of a service. Some should be edited and some should not. The following were prepared by Ignite CLI for your checkers game:


Additional helper functions

Your stored game's black and red fields are only strings, but they represent sdk.AccAddress or even a game from the rules file. Therefore, add helper functions to StoredGame to facilitate operations on them. Create a new file x/checkers/types/full_game.go.

    Get the game's black player:

Copy func (storedGame StoredGame) GetBlackAddress() (black sdk.AccAddress, err error) {
    black, errBlack := sdk.AccAddressFromBech32(storedGame.Black)
    return black, sdkerrors.Wrapf(errBlack, ErrInvalidBlack.Error(), storedGame.Black)
}

- Note how it introduces a new error `ErrInvalidBlack`, which you define shortly. Do the same for the red player. 

2. Parse the game so that it can be played. The `Turn` has to be set by hand:

Note how it introduces a new error ErrInvalidBlack, which you define shortly. Do the same for the red (opens new window) player.

Parse the game so that it can be played. The Turn has to be set by hand:

```go

func (storedGame StoredGame) ParseGame() (game *rules.Game, err errors) {
    board, errBoard := rules.Parse(storedGame.Board)
    if errBoard != nil {
        return nil, sdkerrors.Wrapf(errBoard, ErrGameNotParsable.Error())
    }
    board.Turn = rules.StringPieces[storedGame.Turn].Player
    if board.Turn.Color == "" {
        reutrn nil, sdkerrors.Wrapf(errors.New(fmt.Sprintf("Turn: %s", storedGame.Turn)), ErrGameNotParsable.Error())
    }
    return board, nil
}
```
Add a function that checks a game's validity:
Copy func (storedGame StoredGame) Validate() (err error) {
    _, err = storedGame.GetBlackAddress()
    if err != nil {
        return err
    }
    _, err = storedGame.GetRedAddress()
    if err != nil {
        return err
    }
    _, err = storedGame.ParseGame()
    return err
}

Introduce your own errors:
Copy var (
    ErrInvalidBlack     = sdkerrors.Register(ModuleName, 1100, "black address is invalid: %s")
    ErrInvalidRed       = sdkerrors.Register(ModuleName, 1101, "red address is invalid: %s")
    ErrGameNotParseable = sdkerrors.Register(ModuleName, 1102, "game cannot be parsed")
)


Okay so far as i know what ive done is change the genesis state so that it starts at the next game.

Changes everything for that to work

added some helper functions to Get the black and red addresses, parse the game and validate the game. 



Unit Tests


Now that you have added some code on top of what Ignite CLI created for you, you should add unit tests. You will not add code to test the code generated by Ignite CLI, as your project is not yet ready to test the framework. However, Ignite CLI added some unit tests of its own. Run those for the keeper:

`docker exec -it checkers go test github.com/alice/checkers/x/checkers/keeper`

- It should pass something like 

`ok github.com/alice/checkers/x/checkers/keeper 0.104s`

**Your first unit test** 

- A good start is to test that the default genesis is created as expected. Ignite already created a unit test for the genesis in `x/checkers/types/genesis_test.go`. It runs simple validity test on different genesis examples. 

Table-driven tests basics

- Before digging into the details, let's first discuss a common way of writing tests in Go. A series of related checks can be implemented by looping over a slice of test cases.

```
func TestTime(t *testing.T) {
    testCases := []struct {
        gmt string
        loc string
        want string
    } {
        {"12:31","Europe/Zuri","13:31"}, // incorrect location name
        {"12:31","America/New_York","7:31"}, // should be 07:31
        {"08:08","Australia/Sydney","18:08"},
    }
    for _,tc := range testCases {
        loc, err := time.LoadLocation(tc.loc)
        if err != nil {
            t.Fatalf("could not load location %q", tc.loc)
        }
        gmt, _ := time.Parse("15:04",tc.gmt)
        if got := gmt.In(loc).Format("15:04"); got != tc.want {
            t.Errorf("In(%s, %s) = %s; want %s", tc.gmt, tc.loc, got, tc.want)
        }
    }
}
```
So what this means is essentially there is a struct with gmt loc and want
Then the array of the struct is the europe zuri etc 
Then for the range of test cases 

For each range loop

- Looping over elements in slices, arrays, maps, channels or strings is often better done with a range loop.

```
strings := []string{"hello","world"}
for i, s := range strings {
    fmt.Println(i,s)
}
```

```
0 hello
1 world
```

So the first part of the for loop the is the index.

So in this case

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				SystemInfo: &types.SystemInfo{
					NextId: 39,
				},
				StoredGameList: []types.StoredGame{
					{
						Index: "0",
					},
					{
						Index: "1",
					},
				},
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated storedGame",
			genState: &types.GenesisState{
				StoredGameList: []types.StoredGame{
					{
						Index: "0",
					},
					{
						Index: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

SO this is just where it tests to see if each of the struct states are valid or not.


- You want your tests to pass when everything is okay, but you also want them to fail when something is wrong. Make sure your new test fails by temporarily changing uint64(1) to uint64(2). You should get the following:

```
// The unit test you add is more modest. Your test checks that the starting id on a default
// genesis is 1:
func TestDefaultGenesisState_ExpectedInitialNextId(t *testing.T) {
	require.EqualValues(t,
		&types.GenesisState{
			StoredGameList: []types.StoredGame{},
			SystemInfo:     types.SystemInfo{uint64(1)},
		},
		types.DefaultGenesis())
}
```

- This appears complex, but the significant aspect is this:

```
Diff:
--- Expected
+++ Actual
- NextId: (uint64) 2
+ NextId: (uint64) 1
```

- For *expected* and *actual* to make sense, you hav to ensure that they are correctly placed in your call. When in doubt, go to the `require` function definition:

`func EqualValues(t TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) {...}`

**Debug your unit test**

- Your first unit test is a standard Go unit test. If you use an IDE like Visual Studio Coed 

**More unit tests**

- With a simple yet successful unit test, you can add more concequential ones to test your helper methods. 

- First, create a file that declares some constants that you will reuse throughout:

```
package testutil

const (
    Alice = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3",
    Bob = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
    Carol = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
)
```

- Create a new file `x/checkers/types/full_game_test.go` and declare it in package types_test. Since you are going to repeat some actions, it is worth adding a reusable function:


```
const (
    alice = testutil.Alice
    bob = testutil.Bob
)

func GetStoredGame1() *types.StoredGame {
    return types.StoredGame{
        Black: alice, 
        Red: bob, 
        Index: "1",
        Board: rules.New().String(),
        Turn: "b",
    }
}
```

- Now you can test the function to get the black player's address. One test for the happy path, and another for the error.

https://tutorials.cosmos.network/hands-on-exercise/1-ignite-cli/3-stored-game.html#more-unit-tests

Then what you can do 

ignite chain serve --reset-once

Checks the values saved in SystemInfo

checkersd query checkers --help

show-system-info shows systemInfo

checkersd query checkers show-system-info

```
SystemInfo:
    nextId: "1"
```

The --ouput flag allows you to get your results in a JSON format

checkersd query checkers show-system-info --help

-o, --output string   Output format (text|json) (default "text")

- Now try again a bit differently:

checkersd query checkers show-system-info --output json

{"SystemInfo":{"nextId":"1"}}

checkersd query checkers list-stored-gam

```
pagination:
    next_key: null
    total: "0"
storedGame: []
```

checkersd tx checkers --help

**Create Custom Messages** 


In this section you will:
    - Create a game Protobuf object.
    - Create a game Protobuf srvice interface.
    - Extend your unit tests
    - Interact via the CLI

- You have created your game object type and have decided how to lay games in storage. time to make it possibl for pariticpants to create games.

**Some Initial Thoughts** 

- Because this operation changes the state, it has to originate from transactions an dmessages. Yourm odule recieves a message to create agame- what should go into this message? Questions that you have to answer include:
    - Who is allowed to create a game?
    - Are there any limitations to creating games?
    - Given that a game involves two players, how do you prevent coercion and generally foster good behaviour?
    - Do you want to establish leauges?

- Your implementation does not have to answer everything immediately, but you should be careful that decisions made now do not impede your own future plans or make things more complicated later.

- Keep it simple: a single message should be enough to create a game.

 Code needs

As before:

    What Ignite CLI commands will create your message?
    How do you adjust what Ignite CLI created for you?
    How would you unit-test your addition?
    How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get?

Run the commands, make the adjustments, run some tests. Create the message only, do not create any games in storage for now.


**Create the message** 

- Currently:
    - Your game obejct have ben defined in storage.
    - You prevented a simple CRUD to set the objects straight from transactions.

- Now you need a message to instruct the checkers blockchain to create a game. This message needs to:

    - Not specify the ID of the game, because the system uses an incrementing counter. However, the server needs to return the newly created ID value, since the eventual value cannot be known before hte transaction is included in a block and the state computed. Call this `gameIndex`

    - Not specify the game board as this is controlled by the checkers rules.
    - Specify who is playing with the black pieces. Call the field black.
    - Specify who is playing with the red pieces. Cll the field red.

- Instruct ignite CLI to do all of this 

$ ignite scaffold message createGame black red \
    --module checkers \
    --response gameIndex


**Protobuf Objects**

- Simple Protobuf objects are created:

```
message MsgCreateGame {
    string creator = 1;
    string black = 2;
    string red = 3;
}

message MsgCreateGameResponse {
    string gameIndex = 1;
}

```

- When compiled, for instance with `ignite generate proto-go`, these yield:

```
type MsgCreateGame struct {
     Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
    Black   string `protobuf:"bytes,2,opt,name=black,proto3" json:"black,omitempty"`
    Red     string `protobuf:"bytes,3,opt,name=red,proto3" json:"red,omitempty"`
}

```

and

type MsgCreateGameResponse struct {
    GameIndex string `protobuf:"bytes,1,opt,name=gameIndex,proto3" json:"gameIndex,omitempty"`
}


- Files were generated to serialis the pari whch are named *.pb.go. You should not edit these files.

- Ignite CLI also registered `MsgCreateGame` as a concrete message type with the two (de-)serialization engines:

```
func RegisterCodec(cdc *codec.LegacyAmino) {
    cdc.RegisterConcrete(&MsgCreateGame{}, "checkers/CreateGame",nil)
}
```

and 

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
    registry.RegisterImplementations((*sdk.Msg)(nil),
        &MsgCreateGame{},
    )
    ...
}


This is code that you probably do not need to change.

Ignite CLI also creates boilerplate code to have the message conform to the sdk.Msg (opens new window) type:

```
func (msg *MsgCreateGame) GetSigners() []sdk.AccAddress {
    creator, err := sdk.AccAddressFromBech32(msg.Creator)
    if err != nil {
        panic(err)
    }
    return []sdk.AccAddress{creator}
}
```

- This code is created only once. You can modify it as you see fit. 


**Protobuf Service Interface** 

- Ignite CLI also adds a new function to your gRPC interface that recieves all transcaction messages for the module, because the message is meant to be sent and received. The interface is called service Msg and is declared inside proto/checkers/tx.proto.

Ignite CLI creates this tx.proto (opens new window) file at the beginning when you scaffold your project's module. Ignite CLI separates different concerns into different files so that it knows where to add elements according to instructions received. Ignite CLI adds a function to the empty service Msg with your instruction.

- The new function recieves this `MsgCreateGame`, namely:

```
service Msg {
    rpc CreateGame(MsgCreateGame) returns (MsgCreateGameResponse);
}
```

- As an interface, it does not describe what should happen wehn called. With the help of Protobuf, ignite CLI compiles the interface and creates a default Go implementation.

**Unit tests**

- The code of this section was created by ignite CLI, so there is no point in testing it. However, since you are going to adjust the keeper to do what you want, you should add a test file for that.

- First, recall your address constants in the keeper_test package:

```
package keeper_test

import "github.com/b9lab/checkers/x/checkers/testutil"

const (
    alice = testutil.Alice
    bob   = testutil.Bob
    carol = testutil.Bob
)

```

- Next, create a new `keeper/msg_server_create_game_test.go`, declared with `package keeper_test`:

```
func TestCreateGame(t *testing.T) {
    msgServer, context := setupMsgServer(t)
    createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
        Creator: alice,
        Black:   bob,
        Red:     carol,
    })
    require.Nil(t, err)
    require.EqualValues(t, types.MsgCreateGameResponse{
        GameIndex: "", // TODO: update with a proper value when updated
    }, *createResponse)
}

```

Tested with 

go test github.com/alice/checkers/x/checkers/keeper

- This convenient setupMsgServer function was created by Ignite CLI. To call this a *unit* test is a slight misnomer becaus the `msgServer` created uses a real context and keeper, although with a memory database, not mocks. 


HAve alice start a game with bob

You will have to get the addresses and you can do this with these commands. 

$ export alice=$(docker exec checkers checkersd keys show alice -a)
$ export bob=$(docker exec checkers checkersd keys show bob -a)



It would be to use the crud stuff, then you would edit the code in the keeprs to add your onwn logic and the test are auto created but you can create more in the module somewhere. You then use checkersd to interact with it.





cosmos10chajup2rf4r9e9pm0tstjd28u5swl79wfdm2t -alice

cosmos12z8pkqd56v9swhnc49eznmp6dq5m2yfl087y5n - bob


- How much gas is needed? You can get an estimate by dry running the transaction using the `--dry-run` flag:


checkersd tx checkers create-game $alice $bob --from $alice --dry-run

THis would return the gas estimate


Put the gas on auto to then run the docker command to complete the transaction


docker exec -it checkers checkersd tx checkers create-game $alice $bob --from $alice --gas auto

```
gas estimate: 43032
{"body":{"messages":[{"@type":"/alice.checkers.checkers.MsgCreateGame","creator":"cosmos169mc8qqd6tlued00z23fs75tyecfcazpuwapc4","black":"cosmos169mc8qqd6tlued00z23fs75tyecfcazpuwapc4","red":"cosmos10mqyvj55hm4wunsd62wprwfv9ehcerkfghcjfl"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"43032","payer":"","granter":""}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y
code: 0
codespace: ""
data: 0A280A262F62396C61622E636865636B6572732E636865636B6572732E4D736743726561746547616D65
events:
- attributes:
  - index: true
    key: ZmVl
    value: ""
  type: tx
- attributes:
  - index: true
    key: YWNjX3NlcQ==
    value: Y29zbW9zMTY5bWM4cXFkNnRsdWVkMDB6MjNmczc1dHllY2ZjYXpwdXdhcGM0LzE=
  type: tx
- attributes:
  - index: true
    key: c2lnbmF0dXJl
    value: b1MwcWNrZEtPayt5UlNHdUtNbXZmdFViTjJZbkRTcER0RnNGZVNBais5WWQrQk9vYnRxdHh4Ylp6ZUlib29qd0VNR1BWS1l5Mkg1eHJ3VEZhQ0R5R3c9PQ==
  type: tx
- attributes:
  - index: true
    key: YWN0aW9u
    value: Y3JlYXRlX2dhbWU=
  type: message
gas_used: "41078"
gas_wanted: "43032"
height: "1598"
info: ""
logs:
- events:
  - attributes:
    - key: action
      value: create_game
    type: message
  log: ""
  msg_index: 0
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"create_game"}]}]}]'
timestamp: ""
tx: null
txhash: 576C303E3C43B409B0DEA1CBFF18B7F34F1E69492EE8A562751668117E42834B
```

Returns that 


- If you are curious, the `.events.attrbutes` are encoded in Base64:

```
echo YWN0aW9u | base64 -d
echo YDJFDFKDJFDJKFJK= | base64 -d
```

- Return respectively:

```
action%
create_game%
```

- Which can be found again in `.raw_log`.


- You can query you chain to check whether the system info remains unchanged:

`checkersd query checkers show-system-info`

- This returns:

```
SystemInfo:
    nextid: "1"
```

- It remains unchanged.

- Check whether any game was created:

$ checkersd query checkers list-stored-game

pagination:
  next_key: null
  total: "0"
storedGame: []

- Ahh

- It appears that nothing changed. Ignite CLI created a message, you even signed and broadcast one. However, you have not yet implemented what actions the chain should undertake when it recieves this messag.

- When you are done with this exercise you can stop ignite's `chain serve`


o summarize, this section has explored:

    How to make it possible for participants of the checkers blockchain game to create games with a single message, using a Protobuf object and a Protobuf service interface.
    Which elements must be specified (and which must not) when instructing Ignite CLI to send a game creation message.
    How to add a test file to check the functionality of your code.
    How to interact via the CLI to confirm the "create a game" message occurs as intended - though the absence of a dedicated Message Handler means that currently no game is created.


**Create and Save a Game Properly** 

- In the previous section, you added the message to creat a game along with its serialization and dedicated gRPC function with the help of ignite CLI.

- However, it does not create a game yet because you have not implemented the message hanlding. How would you do this?


**Some initial thoughts** 

- Dwell on the following questions to guide you in this exersise:
    - How do you sanitize your inputs?
    - How do you avoid conflicts with the past and future games?
    - How do you use your files that implement the checkers rules?

**Code Needs** 

- No ignite CLI is involved here, it is just Go.

- Of course, you ned to know where to put your code - look for `TODO`.

- How would you unit-test this message handling?

- How would you use ignite CLI to locally run a one-node blockchain and interact with it via the CLI to se what you get?


- For now, do not bother with niceties like gas metering or event emission.

- You must add code that:
    - Creates a brand new game.
    - Saves it in storage.
    - Returns the ID of the new game.

- Ignite CLI isolated this concern into a seperate file, `x/checkers/keeper/msg_server_create_game.go`, for you to edit:

```
func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)
    // TODO: Handling the message
    _ = ctx
    return &types.MsgCreateGameResponse{}, nil
}
```

- Ignite CLI has conveniently created all the message processing code for you. Yo uare only required to code the key features.

**Coding Steps** 

- Given that you have already done a lot of prepatory work, what coding is involved? How do you replace `//TODO: Handling the message`?

1. First, `rules` represents the ready-made file with the imported ruls of the game:

```
    import (
        ...
+      "github.com/alice/checkers/x/checkers/rules"
        ...
    )
```

2. Get the new game's ID with the Keeper.GetSystemInfo function created by the `ignite scaffold single systeminfo...` command. 

```
systemInfo, found := k.Keeper.GetSystemInfo(ctx)
if !found {
    panic("SystemInfo not found")
}
newIndex := strconv.FormatUint(systemInfo.NextId, 10)
```

- You panic if you cannot find the `SystemInfo` object because there is no way to continue if it is not there. It is not like a user error, which would warrant returning an error.

3. Create the object to be stored:

```
newGame := rules.New()
storedGame := types.StoredGame {
    Index: newIndex, // using the new index from system info here.
    Board: newGame.String(),
    Turn: rules.PieceStrings[newGame.Turn],
    Black: msg.Black,
    Red: msg.Red
}
```

- Note the use of:
    - The rules.New() command, which is part of the checkers rules file you imported earlier.
    - The string content of the msg *types.MsgCreateGame, namely .Black and .Red.

- Also note that you lose the information about the creator. If your design is different, you may want to keep this information.


4. Confirm that the values in the object are correct by checking the validity of players' addresses:

```
err := storedGame.Validate()
if err != nil {
    return nil, err
}
```

`.Red`, and `.Black` need to be checked because they were copied as **strings**. You do not need to check .Creator because at this stage the message signatures have been verified, and the creator is the signer.

- Note that by returning an error, instead of calling `panic`, players cannot stall your blockchain. They can still spam but at a cost, because they will still pay gas fees up to this point.


5. Save the `StoredGame` object using the `Keeper.SetStoredGame` function created by `ignite scaffold map storedGame ...` command: 

`k.Keeper.SetStoredGame(ctx, storedGame)`

6. Prepare the ground for the next game using `Keeper.SetSystemInfo` function created by Ignite CLI:

```
systemInfo.NextId++
k.Keeper.SetSystemInfo(ctx, systemInfo)
```

7. Return the newly created ID for reference:

```
return &types.MsgCreateGameResponse{
    GameIndex: newIndex,
}, nil
```

- You just handled the create game message by actually creating the game. 

**Unit Tests**

- Try the unit tests you prepared in the previous section again: 

`go test github.com/alice/checkers/x/checkers/keeper`

- This should fail with 

```
panic: SystemInfo not found [recovered]
    panic: SystemInfo not found
```

- Your keeper was initialized with an empty genesis. You must fix that one way or another.

- You can fix this by always initialzing the keeper with the default genesis. However such a default initialization may not always be desireable. So it is better to keep this default initialization closest to the tests. Copy the `setupMsgServer` from `msg_server_test.go` into your `msg_server_create_game_test.go`. Modify it to also return the keeper:


```
func setupMsgServerCreateGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
    k, ctx := keepertest.CheckersKeeper(t)
    checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
    return keeper.NewMsgServerImpl(*k) *k, sdk.WrapSDKContext(ctx)
}
```

Note the new import

`"github.com/alice/checkers/x/checkers"`

- Do not forget to replace `setupMsgServer(t)` with this new function everywhere in the file. For instance: 

```
msgServer, _, context := setupMsgServerCreateGame(t)
```


- Run the tests again with the same command as before;


go test github.com/alice/checkers/x/checkers/keeper


- The error has changed to `Not Equal` and you need to adjust the expected value as per the default genesis. 

- The context stores the current context of the application next transaction etc. current execuution.


- One unit test is good, but you can add more, in particular testing whether the values in storage are as expected when you create a single game:

```
func TestCreate1GameHasSaved(t *testing.T) {
    msgSrvr, keeper, context := setupMsgServerCreateGame(t)
    msgSrvr.CreateGame(context, &types.MsgCreateGame{
        Creator: alice,
        Black: bob,
        Red: carol
    })
    systemInfo, found := keeper.GetSystemInfo(sdk.UnwrapSDKContext(context))
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId: 2
    }, systemInfo)
    game1, found1 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
    require.True(t, found1)
    require.EqualValues(t, types.StoredGame{
        Index: "1",
        Board: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
        Turn: "b",
        Black: bob,
        Red: carol,
    }, game1)
}
```

- Or when you create 3 games. Other tests could include whether to *get all* functionality works as expected after you have created 1 game, or 3, or if you create a game in a hypothetical far future. Also add games with badly formatted or missing input.


- `func TestNextToken(t *testing.T)` says the function takes a pointer to the type T in the testing package. T for tests, F for fuzzing, B for benchmarking, etc. That reference is in the variable `t`. 

- `testing.T` stores the state of the test. When Go calls your test functions, it passes the same `testing.T` to each function (presumably). You call methods on it like `t.Fail` to say the test failed, or `t.Skip` to say the test was skipped, etc. It remembers all of this, and Go uses it to report what happend in all the test functions.


- A first good sign is that the output `gas_used` is slightly higher than it was before {`gas_used: "53498"}. After the transaction has been validated, confirm the current state.

- Show the system info:

`checkersd query checkers show-system-info`

this returns 

```
SystemInfo:
    nextId: "2"
```

- List all stored games:

```
checkersd query checkers list-stored-game
```

- This returns a game index `1` as expected:

```
pagination:
  next_key: null
  total: "0"
storedGame:
- black: cosmos169mc8qqd6tlued00z23fs75tyecfcazpuwapc4
  board: '*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*'
  index: "1"
  red: cosmos10mqyvj55hm4wunsd62wprwfv9ehcerkfghcjfl
  turn: b

```

- SHow the new game alone

`checkersd query checkers show-stored-game 1`

- This returns:

```
storedGame:
  black: cosmos169mc8qqd6tlued00z23fs75tyecfcazpuwapc4
  board: '*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*'
  index: "1"
  red: cosmos10mqyvj55hm4wunsd62wprwfv9ehcerkfghcjfl
  turn: b
```

- Now your game is in the blockchain storage. Notice how `alice` was given the black pieces and it is already her turn to play. As a note for the next sections, this is how to understand the board:

```
*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*
                   ^X:1,Y:2                              ^X:3,Y:6
```

- Or if placed in a square:

```


```
X 01234567
*b*b*b*b 0
b*b*b*b* 1
*b*b*b*b 2
******** 3
******** 4
r*r*r*r* 5
*r*r*r*r 6
r*r*r*r* 7
        Y
```

- You can also get this in a one-liner: 

```
docker exec -it checkers \ bash -c "checkersd query checkers show-stored-game 1 --output json 
| jq \".storedGame.board\" | sed 's/\"//g' | sed 's/|/\n/g'"
```

- When you are done with this exersise you can stop Ignite's chain serve. 

synopsis

To summarize, this section has explored:

    How to implement a Message Handler that will create a new game, save it in storage, and return its ID on receiving the appropriate prompt message.
    How to create unit tests to demonstrate the validity of your code.
    How to interact via the CLI to confirm that sending the appropriate transaction will successfully create a game.

**Add a Way to Make a Move** 

- Make sure you have all you need before proceeding:
    - You understand the concepts of transactions, messages, and Protobuf.
    - Go is installed
    - You have the checkers blockchain codebase with MsgCreateGame and its handling. If not, follow the previous stage or check out the relavent version.

- In this section, you will:
    - Extend message handling - play the game.
    - Handle moves and update the game state.
    - Validate input.
    - Extend unit tests.


Your blockchain can now create games, but can you play them? Not yet...so what do you need to make this possible?

**Some Initial Thoughts**

- Before diving into the exerise, take some time to think about the following questions: 
    - What goes into the message?
    - How do you sanitize the inputs?
    - How do you unequivocally identify games?
    - How do you report back errors?
    - How do you use your files that implement the checkers rules?
    - How do you make sure that nothing is lost?

**Code Needs** 

#
Code needs

When it comes to the code you need, ask yourself:

    What Ignite CLI commands will create your message?
    How do you adjust what Ignite CLI created for you?
    How would you unit-test these new elements?
    How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get?

As before, do not bother yet with niceties like gas metering or event emission.

- To play a game a player only needs to specify:
    - The ID of the game the player wants to join. Call the field `gameIndex`.
    - The initial positions of the pawn. Call the fields `fromX` and `fromY` and make them `uint`.
    - The final position of the pawn after a player's move. Call the fields `toX` and `toY` to be `uint` too.

- The player does not need to be explicitly added as a field in the message because the player *is* implicitly the signer of the message. Name the object `PlayMove`.

- Unlike when creating the game, you want to return:
    - The captured piece, if any. Call the fields `capturedX` and `capturedY`.
    - Make then `int` so that you can pass `-1` when no pieces have been captured.
    - The (potential) winner in the field `winner`. 


**With Ignite CLI** 

- Ignite CLI can create the message and the response objects with a single command: 

```
ignite scaffold message playMove gameIndex fromX:uint fromY:uint toX:uint toY:uint \
 --module checkers \
 --response capturedX:int, capturedY:int, winner 
```

- Ignite CLI once more creates all the necessary Protobuf files and boilerplate for you. See `tx.proto`:

```
message MsgPlayMove {
  string creator = 1;
  string gameIndex = 2;
  uint64 fromX = 3;
  uint64 fromY = 4;
  uint64 toX = 5;
  uint64 toY = 6;
}

message MsgPlayMoveResponse {
  int32 capturedX = 1;
  int32 capturedY = 2;
  string winner = 3;
}
```

- All you have to do is fill in the needed part in `x/checkers/keeper/msg_server_play_move.go`

// Youd also have to check .

```
func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

    return &types.MsgPlayMoveResponse{}, nil
}
```

- Where the `TODO` is replaced as per the following


**The Move Handling** 

- The `rules` represent the ready-made file containing the rules of the game you imported earlier.
- Declare your new errors in `x/checkers/types/errors.go`, given your code has to handle new error situations. 

```
var (
    ErrGameNotFound = sdkerrors.Register(ModuleName, 1103, "game by id not found")
    ErrCreatorNotPlayer = sdkerrors.Register(ModuleName, 1104, "message creator not a player")
    ErrNotPlayerTurn = sdkerrors.Register(ModuleName, 1105, "player tried to play out of turn")
    ErrWrongMove = sdkerrors.Register(ModuleName, 1106, "wrong move")
)
```

- Take the following steps to replace the `TODO`: 

1. Fetch the stored game information using the `Keeper.GetStoredGame` function created by Ignite CLI.

Take the following steps to replace the TODO:

    Fetch the stored game information using the Keeper.GetStoredGame (opens new window) function created by Ignite CLI:

Copy storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
if !found {
    return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
}
x checkers keeper msg_server_play_move.go
View source

You return an error because this is a player mistake.

Is the player legitimate? Check with:
Copy isBlack := storedGame.Black == msg.Creator
isRed := storedGame.Red == msg.Creator
var player rules.Player
if !isBlack && !isRed {
    return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
} else if isBlack && isRed {
    player = rules.StringPieces[storedGame.Turn].Player
} else if isBlack {
    player = rules.BLACK_PLAYER
} else {
    player = rules.RED_PLAYER
}
x checkers keeper msg_server_play_move.go
View source

This uses the certainty that the MsgPlayMove.Creator has been verified by its signature (opens new window).

Instantiate the board in order to implement the rules:
Copy game, err := storedGame.ParseGame()
if err != nil {
    panic(err.Error())
}
x checkers keeper msg_server_play_move.go
View source

Fortunately you previously created this helper (opens new window). Here you panic because if the game cannot be parsed the cause may be database corruption.

Is it the player's turn? Check using the rules file's own TurnIs (opens new window) function:
Copy if !game.TurnIs(player) {
    return nil, sdkerrors.Wrapf(types.ErrNotPlayerTurn, "%s", player)
}
x checkers keeper msg_server_play_move.go
View source

Properly conduct the move, using the rules' Move (opens new window) function:
Copy captured, moveErr := game.Move(
    rules.Pos{
        X: int(msg.FromX),
        Y: int(msg.FromY),
    },
    rules.Pos{
        X: int(msg.ToX),
        Y: int(msg.ToY),
    },
)
if moveErr != nil {
    return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
}
x checkers keeper msg_server_play_move.go
View source

Prepare the updated board to be stored and store the information:
Copy storedGame.Board = game.String()
storedGame.Turn = rules.PieceStrings[game.Turn]
k.Keeper.SetStoredGame(ctx, storedGame)
x checkers keeper msg_server_play_move.go
View source

This updates the fields that were modified using the Keeper.SetStoredGame (opens new window) function, as when you created and saved the game.

Return relevant information regarding the move's result:
Copy return &types.MsgPlayMoveResponse{
    CapturedX: int32(captured.X),
    CapturedY: int32(captured.Y),
    Winner:    rules.PieceStrings[game.Winner()],
}, nil
x checkers keeper msg_server_play_move.go
View source

The Captured and Winner information would be lost if you did not get it out of the function one way or another. More accurately, one would have to replay the transaction to discover the values. It is best to make this information easily accessible.


- This completes the move process, facilitated by good preperation and the use of Ignite CLI.

**Unit tests** 

- Adding unit tests for the play message is very similar to what you did for the previous message:
    - Create a new `msg_server_play_move_test.go` file and declare it as package keeper_test. Start with a function that conveniently sets up the keeper for the tests. In this case, already having a game saved can reduce several lines of code in each test:


```
func setupMsgServerWithOneGameForPlayMove(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
    k, ctx := keepertest.CheckersKeeper(t)
    checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
    server := keeper.NewMsgServerImpl(*k)
    context := sdk.WrapSDKContext(ctx)
    server.CreateGame(context, &types.MsgCreateGame {
        Creator: alice,
        Black: bob,
        Red: carol,
    })
    return server, *k, context
}
```

- Note that it reuses `alice`, `bob` and `carol` found in the file `msg_server_create_game_test.go` of the same package.

- Now test the result of a move. Blacks play first, which according to `setupMsgServerWithOneGameForPlayMove` corresponds to `bob`:


Playing a game 

ignite chain serve 

checkersd tx checkers play-move --help

- This returns:

```
Broadcast message playMove

Usage: 
    checkersd tx checkers play-move [game-index] [from-x] [from-y] [to-x] [to-y] [flags]
```

- So Bob tries

$ checkersd tx checkers play-move 1 0 5 1 4 --from $bob

- Game id, from x, from y, to x, to y 

- After you accept sending the transaction, it should complain with the result including:

```
raw_log: 'failed to execute message; message index: 0: {red}: player tried to play out of turn'
...
txhash: D10BB8A706870F65F19E4DF48FB870E4B7D55AF4232AE0F6897C23466FF7871B
```

- If you did not get this raw_log, your transaction may have been sent asynchronously. You can always query a transaction by using the txhash with the following command: 

$ checkersd query tx D10BB8A706870F65F19E4DF48FB870E4B7D55AF4232AE0F6897C23466FF7871B

- And you are back on track 

...
raw_log: 'failed to execute message; message index: 0: {red}: player tried to play
  out of turn'

- Can Alice, who plays *black*, make a move? Can she make a wrong move? For instance, a move from 0-1, to 1-0, which is occupied by one of her pieces. 

$ checkersd tx checkers play-move 1 1 0 0 1 --from $alice

- The computer says "no":

```
...
raw_log: 'failed to execute message; message index: 0: Already piece at destination
  position: {0 1}: wrong move'
```

- So far all seems to be working.

- Time for Alice to make a correct move:

`checkersd tx checkers play-move 1 1 2 2 3 --from alice 

- This returns:

...
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"play_move"}]}]}]'

- Confirm the move went through with your one-line formatter from the previous-section.

$ checkersd query checkers show-stored-game 1 --output json | jq ".storedGame.board" | sed 's/"//g' | sed 's/|/\n/g'

bob's piece moved down and right 

- When you are done with this exersise you can stop ignite's chain serve. 


synopsis

To summarize, this section has explored:

    How to use messages and handlers, in this case to add the capability of actually playing moves on checkers games created in your application.
    The information that needs to be specified for a game move message to function, which are the game ID, the initial positions of the pawn to be moved, and the final positions of the pawn at the end of the move.
    The information necessary to return, which includes the game ID, the location of any captured piece, and the registration of a winner should the game be won as a result of the move.
    How to modify the response object created by Ignite CLI to add additional fields.
    How to implement and check the steps required by move handling, including the declaration of the ready-made rules in the errors.go file so your code can handle new error situations.
    How to add unit tests to check the functionality of your code.
    How to interact via the CLI to confirm that correct player turn order is enforced by the application.


Emit Game Information

- Make sure you have everything you need before proceeeding: 
    - You understand the concept of events.
    - Go is installed 
    - You have the checkers blockchain codebase with MsgPlayMove and its handling. If not, follow the previous steps or check out the relevant version.

- In this section, you will:
    - Define event types.
    - Emit events.
    - Extend unit tests.

- Now that you have added the possible actions, including their return values, use events to notify players. Your blockchain can now create and play games. However, it does not inform the outside about this in a convenient way. 

- This is where events come in - but what do you need to emit them?

- Imagine a potential or current player waiting for their turn. It's not practical to look at all the transactions and search for the ones signifying the player's turn. It is better to listen to known events that let clients determine which player's turn it is. 


Adding events to your application is as simple as:

  1. Defining the events you want to use.
  2. Emitting corresponding events as actions unfold.


**Some Initial Thoughts**

- Before you dive into the specifics of the exercise, ask yourself:
    - Why do actions warrant a detailed event?
    - What level of detail goes into each event?
    - How do you make it easy for external parties to understand your events?
    - At what stage do you emit events?

**Code Needs** 

- Now by thinking about the following: 
    - How do you adjust your code to do all this?
    - How would you unit-test these new elements?
    - How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get? 

- Only focus on the narrow issue of the event emission.


**Game-Created event**

- Start with the event that announces the creation of a new game. The goal is to:
    - Inform the players about the game.
    - Make it easy for the players to find the relevant game.

- Define new keys in `x/checkers/types/keys.go`:

```
const (
    GameCreatedEventType = "new-game-created" // Indicates what event type to listen to
    GameCreatedEventCreator = "creator" // Subsidiary information
    GameCreatedEventGameIndex = "game-index" // What game is relevant
    GameCreatedEventBlock = "black" // Is it relevant to me?
    GameCreatedEventRed = "red" // is it relevant to me?
)
```

- Emit the event in your handler file `x/checkers/keeper/msg_server_create_game.go`:

```
ctx.EventManager().EmitEvent(
    sdk.NewEvent(types.GameCreatedEventType,
        sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
        sdk.NewAttribute(types.GameCreatedEventGameIndex, newIndex),
        sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
        sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
    ),
)
```

- Now you must implement this correspondingly in the GUI, or include a server to listen for such events. 


**Player-moved event** 

- The created transaction to play a move informs us the opponent about: 
    - Which player is relevant.
    - Which game the move relates to.
    - When the move happend.
    - The move's outcome. 
    - Whether the game was won.

- Contrary to the *create game* event, which alerted the players about the new game, the players now know which game IDs to watch for. There is no need to repeat the player's addresses, the game ID is information enough.

- You define new keys in `x/checkers/types/keys.go` similarly: 

```
const (
    MovePlayedEventType      = "move-played"
    MovePlayedEventCreator   = "creator"
    MovePlayedEventGameIndex = "game-index"
    MovePlayedEventCapturedX = "captured-x"
    MovePlayedEventCapturedY = "captured-y"
    MovePlayedEventWinner    = "winner"
)
```

- Emit the event in your life `x/checkers/keeper/msg_server_play_move.go`.

```
ctx.EventManager().EmitEvent(
    sdk.NewEvent(types.MovePlayedEventType,
        sdk.NewAttribute(types.MovePlayedEventCreator, msg.Creator),
        sdk.NewAttribute(types.MovePlayedEventGameIndex, msg.GameIndex),
        sdk.NewAttribute(types.MovePlayedEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
        sdk.NewAttribute(types.MovePlayedEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
        sdk.NewAttribute(types.MovePlayedEventWinner, rules.PieceStrings[game.Winner()]),
    ),
)
```

**Unit Tests** 


- The unit tests you have created so far still pass. However you also want to confirm that the events have been emitted in both situations. The events are recoreded in the context, so the test is a little bit different. In `msg_server_create_game_test.go`, add this test:

```
func TestCreate1GameEmitted(t *testing.T) {
    msgSrvr, _, context := setupMsgServerCreateGame(t)
    msgSrvr.CreateGame(context, &types.MsgCreateGame{
        Creator: alice,
        Black:   bob,
        Red:     carol,
    })
    ctx := sdk.UnwrapSDKContext(context)
    require.NotNil(t, ctx)
    events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
    require.Len(t, events, 1)
    event := events[0]
    require.EqualValues(t, sdk.StringEvent{
        Type: "new-game-created",
        Attributes: []sdk.Attribute{
            {Key: "creator", Value: alice},
            {Key: "game-index", Value: "1"},
            {Key: "black", Value: bob},
            {Key: "red", Value: carol},
        },
    }, event)
}


```


- How can you *guess* the order of the elements? Easily, as you created them in this order. Alternatively, you can *peek* using Visual Studio Code:
    1. Put a break point on the line after `event := events[0]`
    2. Run this test in **debug mode**: right-click the green arrow next to the test name.

As for the events emitted during the play move test, there are two of them: one for the creation and the other for the play. Because this is a unit test and each action is not isolated into individual transactions, the context collects all events emitted during the test. It just so happens that the context prepends them - the newest one is at index 0. Which is why, when you fetch them, the play event is at events[0].

func TestPlayMoveEmitted(t *testing.T) {
    msgServer, _, context := setupMsgServerWithOneGameForPlayMove(t)
    msgServer.PlayMove(context, &types.MsgPlayMove{
        Creator:   bob,
        GameIndex: "1",
        FromX:     1,
        FromY:     2,
        ToX:       2,
        ToY:       3,
    })
    ctx := sdk.UnwrapSDKContext(context)
    require.NotNil(t, ctx)
    events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
    require.Len(t, events, 2)
    event := events[0]
    require.EqualValues(t, sdk.StringEvent{
        Type: "move-played",
        Attributes: []sdk.Attribute{
            {Key: "creator", Value: bob},
            {Key: "game-index", Value: "1"},
            {Key: "captured-x", Value: "-1"},
            {Key: "captured-y", Value: "-1"},
            {Key: "winner", Value: "*"},
        },
    }, event)
}

When two players play one after the other, the context collates the attributes of move-played all together in a single array in an appending fashion, with the older attributes at the lower indices, starting at 0. For instance, you have to rely on array slices like event.Attributes[5:] to test the attributes of the second move-played event:


func TestPlayMove2Emitted(t *testing.T) {
    msgServer, _, context := setupMsgServerWithOneGameForPlayMove(t)
    msgServer.PlayMove(context, &types.MsgPlayMove{
        Creator:   bob,
        GameIndex: "1",
        FromX:     1,
        FromY:     2,
        ToX:       2,
        ToY:       3,
    })
    msgServer.PlayMove(context, &types.MsgPlayMove{
        Creator:   carol,
        GameIndex: "1",
        FromX:     0,
        FromY:     5,
        ToX:       1,
        ToY:       4,
    })
    ctx := sdk.UnwrapSDKContext(context)
    require.NotNil(t, ctx)
    events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
    require.Len(t, events, 2)
    event := events[0]
    require.Equal(t, "move-played", event.Type)
    require.EqualValues(t, []sdk.Attribute{
        {Key: "creator", Value: carol},
        {Key: "game-index", Value: "1"},
        {Key: "captured-x", Value: "-1"},
        {Key: "captured-y", Value: "-1"},
        {Key: "winner", Value: "*"},
    }, event.Attributes[5:])
}

- Try these tests:

go test github.com/alice/checkers/x/checkers/keeper


**Interact with the CLI** 


- If you did not do it already, start your chain with ignite

Alice made a move. Will Bob's move emit an event?

Copy $ checkersd tx checkers play-move 1 0 5 1 4 --from $bob

The log is longer and not very readable, but the expected elements are present:
Copy ...
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"play_move"}]},{"type":"move-played","attributes":[{"key":"creator","value":"cosmos1xf6s64kaw7at7um8lnwj65vadxqr6hnyhr9v83"},{"key":"game-index","value":"1"},{"key":"captured-x","value":"-1"},{"key":"captured-y","value":"-1"},{"key":"winner","value":"*"}]}]}]'

To parse the events and display them in a more user-friendly way, take the txhash again:
Copy $ checkersd query tx 531E5708A1EFBE08D14ABF947FBC888BFC69CD6F04A589D478204BF3BA891AB7 --output json | jq ".raw_log | fromjson"

Copy $ docker exec -it checkers \
    bash -c "checkersd query tx 531E5708A1EFBE08D14ABF947FBC888BFC69CD6F04A589D478204BF3BA891AB7 --output json | jq '.raw_log | fromjson'"

This returns something like:
Copy [
  {
    "events": [
      {
        "type": "message",
        "attributes": [
          {
            "key": "action",
            "value": "play_move"
          }
        ]
      },
      {
        "type": "move-played",
        "attributes": [
          {
            "key": "creator",
            "value": "cosmos1xf6s64kaw7at7um8lnwj65vadxqr6hnyhr9v83"
          },
          {
            "key": "game-index",
            "value": "1"
          },
          {
            "key": "captured-x",
            "value": "-1"
          },
          {
            "key": "captured-y",
            "value": "-1"
          },
          {
            "key": "winner",
            "value": "*"
          }
        ]
      }
    ]
  }
]

As you can see, no pieces were captured. However, it turns out that Bob placed his piece ready to be captured by Alice:

$ checkersd query checkers show-stored-game 1 --output json | jq ".storedGame.board" | sed 's/"//g' | sed 's/|/\n/g'

Which prints

*b*b*b*b
b*b*b*b*
***b*b*b
**b*****
*r******    <-- Ready to be captured
**r*r*r*
*r*r*r*r
r*r*r*r*

storedGame:
  black: cosmos10chajup2rf4r9e9pm0tstjd28u5swl79wfdm2t
  board: '*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*'
  index: "0"
  red: cosmos12z8pkqd56v9swhnc49eznmp6dq5m2yfl087y5n
  turn: b

- The rules in this game included in this project mandate that the player captures a piece when possible. So Alice captures the piece:

`checkersd tx checkers play-move 1 2 3 0 5 --from $alice`

```
...
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"play_move"}]},{"type":"move-played","attributes":[{"key":"creator","value":"cosmos1qxeu0aclpl45429aeveh3t4e7y9ghr22r5d9r2"},{"key":"game-index","value":"1"},{"key":"captured-x","value":"1"},{"key":"captured-y","value":"4"},{"key":"winner","value":"*"}]}]}]'
```

docker exec -it checkers bash -c "checkersd query tx sjkfsdjk --output json | jq '.raw_log | fromjson'

When formatted 
[
  {
    "events": [
      {
        "type": "message",
        "attributes": [
          {
            "key": "action",
            "value": "play_move"
          }
        ]
      },
      {
        "type": "move-played",
        "attributes": [
          {
            "key": "creator",
            "value": "cosmos10chajup2rf4r9e9pm0tstjd28u5swl79wfdm2t"
          },
          {
            "key": "game-index",
            "value": "0"
          },
          {
            "key": "captured-x",
            "value": "1"
          },
          {
            "key": "captured-y",
            "value": "4"
          },
          {
            "key": "winner",
            "value": "*"
          }
        ]
      }
    ]
  }
]

- Correct: Alice captured a piece and the board now looks like this: 


*b*b*b*b
b*b*b*b*
***b*b*b
********
********
b*r*r*r*
*r*r*r*r
r*r*r*r*

- This confirms that the *play* event is emitted as expected. You can confirm the same for the *game created* event.

- When you are done with this exercise you can stop Ignite's `chain serve`


synopsis

To summarize, this section has explored:

    - How to define event types and then emit events to cause the UI to notify players of game actions as they occur, such as creating games and playing moves.
    - How listening to known events which let clients determine which player must move next is better than the impractical alternative of examining all transactions to search for the ones which signify a player's turn.
    - How to define a Game-created event that will notify the participating players and make it easy for them to find the game.
    - How to define a Player-moved event that will indicate which player and game is involved, when the move occurred, the move's outcome, and whether the game was won as a result.
    - How to test your code to ensure that it functions as desired.
    - How to interact with the CLI to check the effectiveness of an emitted event.


**Make Sure a Player Can Reject a Game** 


- Before proceeding, make sure you have all you need:
    - You understand the concepts of transactions, messages, and Protobuf.
    - You know how to create a message with Ignite CLI, and code its handling. This section does not aim to repeat what can be learned in earlier sections.
    - Go is installed 
    - You have the checkers blockchain codebase with the previous messages and their events. If not, follow the previous steps or check out the relevant section.


- In this section you will: 
    - Add a new protocol rule.
    - Define custom errors. 
    - Add a message handler.
    - Extend unit tests.


- Your blockchain can now create and play games, and inform the outside world about the process. It would be good to add a way for players to back out of games they do not want to play. What do you need to make this possible? 

**Some Initial Thoughts** 

- Ask yourself:
    - What goes into the messages? 
    - How can you santize the inputs?
    - How do you unoquivocally identify games?
    - What conditions have to be satisfied to reject a game?
    - How do you report back errors?
    - What event should you emit?
    - How do you use your files that implement the checkers rules?
    - What do you you do a rejected game?

**Code Needs** 

- When you think about the code you might need, try to first answer the following questions:
    - What Ignite CLI commands will create your messages?
    - How do you adjust what Ignite CLI created for you?
    - How would you unit-test these new elements?
    - How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get?

- As before, do not bother yet with niceties of gas metering?

- If anyone can create a game for any two other players, it is important to allow a player to reject a game. But a player should not be allowed to reject a game once they have made their first move.

- To reject a game, a player needs to provide the ID of the game that the player wants to reject. Call the field `gameIndex`. This should be sufficient, as the signer of the message is implictly the player. 

**Working with Ignite CLI** 

- Name the message object `RejectGame`. Invoke Ignite CLI:

`ignite scaffold message rejectGame gameIndex --module checkers`

- THis creates all the boilerplate for you and leaves a single place for the code you want to include: 

```
func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

    return &types.MsgRejectGameResponse{}, nil
}

```

**Additional Information**

- A new rule of the game should be that the player cannot reject a game once they begin to play. When loading a `StoredGame` from storage you have no way of knowing whether a player already played or not. To access this information add a new field to the `StoredGame` called `MoveCount`. In `proto/checkers/stored_game.proto`.

```
  message StoredGame {
        ...
+      uint64 moveCount = 6;
    }
```

- Run protobuf to recompile the relevant Go files: 

ignite generate proto-go


- `MoveCount` should start at `0` and increment by `1` on each move. 

1. Adjust it first in the handler when creating the game: 

```
storedGame := types.StoredGame {
    ...
    MoveCount: 0,
}
```

2. Before saving to the storage, adjust it in the handler when playing a move: 

```
...
storedGame.MoveCount++
storedGame.Board = game.String()
...
```


- With `MoveCount` counting properly, you are now ready to handle a rejection request. 


**The reject handling** 

- To follow the Cosmos SDK conventions, declare the following new errors: 

```
var (
    ErrBlackAlreadyPlayed = sdkerrors.Register(ModuleName, 1107, "black player has already played")
    ErrRedAlreadyPlayed = sdkerrors.Register(ModuleName, 1108, "red player has already played")
)
```

- This time you will add an event for rejection. Begin by preparing the new keys: 


```
const (
    GameRejectedEventType = "game-rejected"
    GameRejectedEventCreator = "creator"
    GameRejectedEventGameIndex = "game-index"
)

```

- In the message handler, the reject steps are:

1. Fetch the relevant information:
    ```
    storedGame, found := k.Keeper.GetStoredGame(ctx, msg.GameIndex)
    if !found {
        return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "%s", msg.GameIndex)
    }

    ```

2. Is the player expected? Did the player already play? Check with: 

```
if storedGame.Black == msg.Creator {
    if 0 < storedGame.MoveCount { // Notice the use of the new field
        return nil, types.ErrBlackAlreadyPlayed
    }
} else if storedGame.Red == msg.Creator {
    if 1 < storedGame.MoveCount { // Notice the use of the new field
        return nil, types.ErrRedAlreadyPlayed
    }
} else {
    return nil, sdkerrors.Wrapf(types.ErrCreatorNotPlayer, "%s", msg.Creator)
}
```

- Remember that the player with the color black plays first. 

3. Remove the game using the `Keeper.RemoveStoredGame` function created long ago by the `ignite scaffold map storedGame ...` command: 

`k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)`

4. Emit the relevant event:

```
ctx.EventManager().EmitEvent(
    sdk.NewEvent(types.GameRejectedEventType,
        sdk.NewAttribute(types.GameRejectedEventCreator, msg.Creator),
        sdk.NewAttribute(types.GameRejectedEventGameIndex, msg.GameIndex),
        ),
)
```

5. Leave the returned object as it is, as you have nothing new to tell the caller.


- Finally, confirm that your project at least compiles with :

ignite chain build

**Unit Tests** 


- Before testing what you did when rejecting a game, you have to fix the existing tests by adding `MoveCount: 0`, or more when testing a retrieved `StoredGame`.

- When you are done with the existing tests, the tests for *reject* here are similar to thoes you created for *create and play*, except that now you test a game rejection by the game creator, the black player, or the red player which is made before anyone has played, or after one or two moves have been made. Check also that the game is removed, and that elements are emitted.

- For instance: 

func TestRejectGameByRedOneMoveRemovedGame(t *testing.T) {
    msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
    msgServer.PlayMove(context, &types.MsgPlayMove{
        Creator:   bob,
        GameIndex: "1",
        FromX:     1,
        FromY:     2,
        ToX:       2,
        ToY:       3,
    })
    msgServer.RejectGame(context, &types.MsgRejectGame{
        Creator:   carol,
        GameIndex: "1",
    })
    systemInfo, found := keeper.GetSystemInfo(sdk.UnwrapSDKContext(context))
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId: 2,
    }, systemInfo)
    _, found = keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
    require.False(t, found)
}


- Try these tests:

go test github.com/alice/checkers/x/checkers/keeper


**Interact with the CLI**

- Time to see if it is possible to reject a game from the command line. If you did not do it already, start your chain with ignite.

- First, it is possible to reject the current game from the command line?

`checkersd tx checkers --help`

- This prints: 

...
Available Commands:
...
  reject-game Broadcast message rejectGame

  reject-game is the command. What is its syntax?

Copy $ checkersd tx checkers reject-game --help

This prints:
Copy ...
Usage:
  checkersd tx checkers reject-game [game-index] [flags]

Have Bob, who played poorly in game 1, try to reject it:

Copy $ checkersd tx checkers reject-game 1 --from $bob

This returns:
Copy ...
raw_log: '[{"events":[{"type":"game-rejected","attributes":[{"key":"creator","value":"cosmos14g3qw6nkk8zc762k87cg77w7vd8xdnffnp2w6u"},{"key":"game-index","value":"1"}]},{"type":"message","attributes":[{"key":"action","value":"reject_game"}]}]}]'

Against expectations, the system carried out Bob's request to reject the game. Confirm that the game has indeed been removed from storage:

Copy $ checkersd query checkers show-stored-game 1

This returns:
Copy Error: rpc error: code = NotFound desc = rpc error: code = NotFound desc = not found: key not found
...


- How is it possible that Bob could reject a game he already had played in, despite the code preventing that? Because game 1 was created in an earlier version of your code. 

- This earlier version created **a game without any .MoveCount,** or more precisely with `MoveCount == 0`. When you later added the code for rejection, Ignite CLI kept the current state of your blockchain. In effect, your blockchain was in a *broken* state, where **the code and the state were out of sync**.
- To see how to properly handle code changes that would otherwise result in a broken state, see the section on migrations.


You have to create other games and test the rejection on them. Notice the incrementing game ID.


bob
cosmos12z8pkqd56v9swhnc49eznmp6dq5m2yfl087y5n

alice
cosmos10chajup2rf4r9e9pm0tstjd28u5swl79wfdm2t

You have to create other games and test the rejection on them. Notice the incrementing game ID.
1

Black rejects:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
$ checkersd tx checkers reject-game 2 --from $alice

Above, Alice creates a game and rejects it immediately. This returns:
Copy ...
raw_log: '[{"events":[{"type":"game-rejected","attributes":[{"key":"creator","value":"cosmos1uhfa4zhsvz7cyec7r62p82swk8c85jaqt2sff5"},{"key":"game-index","value":"2"}]},{"type":"message","attributes":[{"key":"action","value":"reject_game"}]}]}]'

Correct result, because nobody played a move.
2

Red rejects:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
$ checkersd tx checkers reject-game 3 --from $bob

Above, Alice creates a game and Bob rejects it immediately. This returns:
Copy ...
raw_log: '[{"events":[{"type":"game-rejected","attributes":[{"key":"creator","value":"cosmos14g3qw6nkk8zc762k87cg77w7vd8xdnffnp2w6u"},{"key":"game-index","value":"3"}]},{"type":"message","attributes":[{"key":"action","value":"reject_game"}]}]}]'

Correct again, because nobody played a move.
3

Black plays and rejects:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
$ checkersd tx checkers play-move 4 1 2 2 3 --from $alice
$ checkersd tx checkers reject-game 4 --from $alice

Above, Alice creates a game, makes a move, and then rejects the game. This returns:
Copy ...
raw_log: 'failed to execute message; message index: 0: black player has already played'

Correct: the request fails, because Alice has already played a move.
4

Alice plays and Bob rejects:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
$ checkersd tx checkers play-move 5 1 2 2 3 --from $alice
$ checkersd tx checkers reject-game 5 --from $bob

Above, Alice creates a game, makes a move, and Bob rejects the game. This returns:
Copy ...
raw_log: '[{"events":[{"type":"game-rejected","attributes":[{"key":"creator","value":"cosmos14g3qw6nkk8zc762k87cg77w7vd8xdnffnp2w6u"},{"key":"game-index","value":"5"}]},{"type":"message","attributes":[{"key":"action","value":"reject_game"}]}]}]'

Correct: Bob has not played a move yet, so he can still reject the game.
5

Alice & Bob play, Bob rejects:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
$ checkersd tx checkers play-move 6 1 2 2 3 --from $alice
$ checkersd tx checkers play-move 6 0 5 1 4 --from $bob
$ checkersd tx checkers reject-game 6 --from $bob

Above, Alice creates a game and makes a move, then Bob makes a poor move and rejects the game. This returns:
Copy ...
raw_log: 'failed to execute message; message index: 0: red player has already played'

Correct: this time Bob could not reject the game because the state recorded his move in .MoveCount.


- To belabor the point made in the earlier box: if you change your code, think about what it means for the current state of the chain and whether you end up with a broken state.

- In this case, you could first introduce the MoveCount and its handling. Then when all games have been correctly counted, you introduce the rejection mechanism.


synopsis

To summarize, this section has explored:

    How to use messages and handlers to build on the gameplay functionalities of your application by adding the capacity for players to reject participating in a game.
    How to create a new RejectGame message object including ID of the game to be rejected.
    How to add a new rule with the necessary additional information to prevent players from backing out of games in which they have already played moves, and how to declare new errors that respond to attempts to break this new rule.
    How to add a unit test to check that games can be rejected by the game creator, the black player, and the red player under the approved circumstances, and to check that rejected games are removed and that events are emitted.
    How to interact via the CLI to confirm the new "game rejection" function is performing as required, and to be aware that preexisting games will permit incorrect game rejection due to your blockchain being in a broken state due to your subsequent changes.


*CONTINUE DEVELOPING YOUR COSMOS CHAIN*

- You will work further on your checkers blockchain and make your next steps with Ignite CLI. You have a workable checkers blockchain, one which lets players play. 

- But have you thought about everything? Is your blockchain safe from bad behaviour? How do you incentivize good behaviour? Can you also make it more fun? 

- Continue your journey with Ignite CLI: learn how to introduce a wager, manage gas, and query for players' moves. 


- In this chapter 


In this chapter, you will:

    Continuously develop your checkers blockchain with the Ignite CLI.
    Let players set a wager.
    Order your games and introduce a game deadline.
    Record the winners.
    Help players do a correct move.
    Explore how you can manage gas for your application-specific chain.


**Put Your Games In Order** 

- Make sure you have everything you need before proceeding: 
    - You understand the concepts of ABCI, Protobuf, and of a doubly-linked-list
    - Go is installed
    - You have the checkers blockchain codebase with `MsgRejectGame` and its handling.
    - If not follow the previous steps or checkout the relevant version.


- In this section you will deal with:
    - The FIFO data structure
    - FIFO unit tests

- You will learn:
    - Modularity and data orginization styles.

- In the previous step, you added a way for players to reject a game, so there is a way for stale games to be removed from storage. But is this enough to avoid *state pollution*?

- There are some initial thoughts and code needs to keep in mind during the next sections to be able to implement forfeits in the end. 


**Some Initial Thoughts** 

- Before you begin touching your code, ask:
    - What conditions have to be satisfied for a game to be considered stale and the blockchain to act?
    - How do you sanitize the new information inputs? 
    - How would you get rid of stale games as part of the protocol, that is *without user inputs?* 
    - How do you optimize performance and data structures so that a few stale games do not cause your blockchain to grind to a halt?
    - How can you be sure that your blockchain is safe from attacks? 
    - How do you make your changes compatible with future plans for wagers?
    - Are there errors to report back? 
    - What event should you emit? 

**Code Needs** 

- Now, think about what possible code changes and additions you should consider: 
    - What Ignite CLI commands, if any, will assist you? 
    - How do you adjust what Ignite CLI created for you? 
    - How would you unit-test these new elements?
    - How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get?

- For now, do not bother yet with future ideas like wager handling? 


**Why would you reject** 

- There are two ways for a game to advance through its lifecycle until resolution, win or draw: *play and reject*.

- Game inactivity could become a factor. What if a player never shows up again? Should a game remain in limbo forever?
 
- Eventually you want to let players wager on the outcome of games, so you don't want games remaining in limbo if they have *value* assigned. For this reason, you need a way for games to be forcibly resolved if one player stops responding.

- The simplest mechanism to expire a game is to use a **deadline**. If the deadline is reached, then the game is forcibly terminated and expires. The deadline is pushed back every time a move is played.

- To enforce the termination, it is a good idea to use the `**EndBlock**` part of the ABCI protocol. The call **EndBlock** is triggered when all transactions of the block are delivered, and allows you to tidy up before the block is sealed. In your case, all games that have reached their deadline will be terminated. 

- How do you find all the games that have reached their deadline? You could use a pseudo-code like: 

`findAll(game => game.deadline < now)`

- This approach is **expensive** in terms of computation. The `EndBlock` code should not have to pull up all games out of storage just to find a few that are relevant. 

- Doing a `findAll` costs `O(n)`, where `n` is the total number of games.


**How can you reject?** 

- You need another data structure. The simplest option is a First-In-First-Out (FIFO) that is constantly updated, so that: 
    - When games are played, they are taken out of where they are and sent to the tail.
    - Games that have not been played for the longest time eventually rise to the head. 

- Therefore, when terminating expired games in `EndBlock`, you deal with the expired games that are the head of the FIFO. You do not stop until the head includes an ongoing game. The cost is: 
    - O(1) on each game creation and gameplay.
    - O(k) where k is the number of expired games on each block.
    - k =< n where n is the number of games that exist.

- k is still an unbounded number of operations. However, if you use the same expiration duration on each game, for `k` games to expire together in a given block they would all have to have had a move in the same previous block (give or take the block before or after). In the worst case, the largest `EndBlock` computation will be proportional to the largest regular block in the past. This is a reasonable risk to take. 

- This only works if the expiration duration is the same for all games, instead of being a parameter left to a potentially malicious game creator. 

- Well you could make it that there is a limit to how many games expire in a block so you can limit slow chain processing. by tagging ones that are expired that didnt make this block like skipping ones that need to be removed but not in this block???


**New Information** 

- How do you implement a FIFO from which you extract elements at random positions? Choose a doubly-linked list:
    1. You must remember the game ID at the head to pick expired games, and at the tail to send back fresh games. The existing `SystemInfo` object is useful, as it is already expandable. Add to its Protobuf decleration:

```
    message SystemInfo {
        ...
        string fifoHeadIndex = 2; // Will contain the index of the game at the head.
        string fifoTailIndex = 3; // Will contain the index of the game at the tail.
    }
```

2. To make extraction possible, each game must know which other game takes place before it is in the FIFO, and which after. Store this double-link information in `StoredGame`. Add them to the game's Protobuf decleration.


```
    message StoredGame {
        ...
        string beforeIndex = 7; // Pertains to the FIFO. Toward head.
        string afterIndex = 8; // Pertains to the FIFO. Toward tail.
    } 

```

3. There must be an "ID" that indicates *no game*. Uses `"-1"`, which you save as a constant:

```
const (
    NoFifoIndex = "-1"
)
```

4. Instruct Ignite CLI and Protobuf to regenerate the Protobuf files:

`ignite generate proto-go`



Adjust the default genesis values, so that it has a proper head and tail:
Copy     func DefaultGenesis() *GenesisState {
        return &GenesisState{
            SystemInfo: SystemInfo{
                NextId:        uint64(DefaultIndex),
+              FifoHeadIndex: NoFifoIndex,
+              FifoTailIndex: NoFifoIndex,
            },
            ...
        }
    }
x checkers types genesis.go 

**FIFO Management** 

- Now that the new fields are created, you need to update them to keep your FIFO up-to-date. It's better to create a seperate file that encapsulates this knowledge. 

- Create `x/checkers/keeper/stored_game_in_fifo.go` with the following: 

1. A function to remove from the FIFO:

```
func (k Keeper) RemoveFromFifo(ctx sdk.Context, game *types.StoredGame, info *types.SystemInfo) {
    // Does it have a predesessor?
    if game.BeforeIndex != types.NoFifoIndex {
        beforeElement, found := k.GetStoredGame(ctx, game.BeforeIndex)
        if !found {
            panic("Element before in Fifo was not found")
        }
        beforeElement.AfterIndex = game.AfterIndex
        k.SetStoredGame(ctx, beforeElement)
        if game.AfterIndex == types.NoFifoIndex {
            info.FifoTailIndex = beforeElement.Index
        }
        // Is it at the FIFO head?
    } else if info.FifoHeadIndex == game.Index {
      info.FifoHeadIndex = game.AfterIndex  
    }
    // Does it have a successor?
    if game.AfterIndex != types.NoFifoIndex {
        afterElement, found := k.GetStoredGame(ctx, game.AfterIndex)
        if !found {
            panic("Element after in Fifo was not found")
        }
        afterElement.BeforeIndex = game.BeforeIndex
        k.SetStoredGame(ctx, afterElement) 
        if game.BeforeIndex == types.NoFifoIndex {
            info.FifoHeadIndex = afterElement.Index
        }
        // Is it at the FIFO tail?
    } else if info.FifoTailIndex == game.Index {
        info.FifoTailIndex = game.BeforeIndex
    }
    // essentially delete
    game.BeforeIndex = types.NoFifoIndex
    game.AfterIndex = types.NoFifoIndex
}

```

- The game is passed as an argument is **not** saved in storage here, even if it was updated.
- Only its fields in memory are adjusted. The *before* and *after* games are saved in storage.
- Do a `SetStoredGame` after calling this function to avoid having a mix of saves and memory states. The same applies to `SetSystemInfo`.

2. A function to send to the tail:

// So this essentially takes the game out of the doubly linked list then adds it to the tail.
```go 
func (k Keeper) SendToFifoTail(ctx sdk.Context, game *types.StoredGame, info *types.SystemInfo) {
    // Essentially if the head and tail don't exist yet, list is empty.
    if info.FifoHeadIndex == types.NoFifoIndex && info.FifoTailIndex == types.NoFifoIndex {
        game.BeforeIndex = types.NoFifoIndex
        game.AfterIndex = types.NoFifoIndex
        info.FifoHeadIndex = game.Index
        info.FifoTailIndex = game.Index
    // if the head or tail are empty.
    } else if info.FifoHeadIndex == types.NoFifoIndex || info.FifoTailIndex == types.NoFifoIndex {
        panic("Fifo should have both head and tail or none")
    } else if info.FifoTailIndex == game.Index {
        // Nothing to do, already at tail
    } else {
        // Snip game out
        k.RemoveFromFifo(ctx, game, info)

        // Now add to tail
        currentTail, found := k.GetStoredGame(ctx, info.FifoTailIndex)
        if !found {
            panic("Current Fifo tail was not found")
        }
        currentTail.AfterIndex = game.Index
        k.SetStoredGame(ctx, currentTail)

        game.BeforeIndex = currentTail.Index
        info.FifoTailIndex = game.Index
    }
}
```

- Again, it is advisable to do `SetStoredGame` and `SetSystemInfo` after calling this function.


**FIFO integration** 

- With these functions ready, it is time to use them in the message handlers.

    1. In the handler when creating a new game, set default values for `BeforeIndex` and `AfterIndex`. 

```go
...
    storedGame := types.StoredGame {
        ...
        BeforeIndex: types.NoFifoIndex,
        AfterIndex: types.NoFifoIndex,
    }
```


Send the new game to the tail because it is freshly created:
Copy     ...
+  k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)
    k.Keeper.SetStoredGame(ctx, storedGame)
    ...
x checkers keeper msg_server_create_game.go
View source

In the handler, when playing a move send the game back to the tail because it was freshly updated:
Copy     ...
+  systemInfo, found := k.Keeper.GetSystemInfo(ctx)
+  if !found {
+      panic("SystemInfo not found")
+  }
+  k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)

    storedGame.MoveCount++
    ...
    k.Keeper.SetStoredGame(ctx, storedGame)
+  k.Keeper.SetSystemInfo(ctx, systemInfo)
    ...
x checkers keeper msg_server_play_move.go
View source

Note that you also need to call SetSystemInfo.

In the handler, when rejecting a game remove the game from the FIFO:

    Copy     ...
    +  systemInfo, found := k.Keeper.GetSystemInfo(ctx)
    +  if !found {
    +      panic("SystemInfo not found")
    +  }
    +  k.Keeper.RemoveFromFifo(ctx, &storedGame, &systemInfo)
        k.Keeper.RemoveStoredGame(ctx, msg.GameIndex)
    +  k.Keeper.SetSystemInfo(ctx, systemInfo)
        ...
    x checkers keeper msg_server_reject_game.go
    View source

You have implemented a FIFO that is updated but never really used. It will be used in a later section.


**Unit Tests** 

- At this point, your previous unit tests are failing, so they must be fixed. Add `FifoHeadIndex` and `FifoTailIndex` in your value requirements on `SystemInfo` as you create games, play moves, and reject games. 

- Also add `BeforeIndex` and `AfterIndex` in your value requirements on `StoredGame` as you create games and play moves. 



alice
cosmos17ww205xgl708tndsanhyhjyww8s832a2edk8yp
bob
cosmos1dk68clfuky4nmpv45qdasadjj4u63khueyk4rx


https://tutorials.cosmos.network/hands-on-exercise/2-ignite-cli-adv/1-game-fifo.html#


teract via the CLI

Time to explore the commands. You need to start afresh because you made numerous additions to the blockchain state:

Copy $ ignite chain serve --reset-once

Do not forget to export alice and bob again, as explained in an earlier section under "Interact via the CLI".
1

Is the genesis FIFO information correctly saved?

Copy $ checkersd query checkers show-system-info

This should print:
Copy SystemInfo:
    fifoHeadIndex: "-1" # There is nothing
    fifoTailIndex: "-1" # There is nothing
    nextId: "1"
2

If you create a game, is the game as expected?

Copy $ checkersd tx checkers create-game $alice $bob --from $bob
$ checkersd query checkers show-system-info

This should print:
Copy SystemInfo:
    fifoHeadIndex: "1" # The first game you created
    fifoTailIndex: "1" # The first game you created
    nextId: "2"
3

What about the information saved in the game?

Copy $ checkersd query checkers show-stored-game 1

Because it is the only game, this should print:
Copy storedGame:
    afterIndex: "-1" # Nothing because it is alone
    beforeIndex: "-1" # Nothing because it is alone
    ...
4

And if you create another game?

Copy $ checkersd tx checkers create-game $alice $bob --from $bob
$ checkersd query checkers show-system-info

This should print:
Copy SystemInfo:
    fifoHeadIndex: "1" # The first game you created
    fifoTailIndex: "2" # The second game you created
    nextId: "3"
5

Did the games also store the correct values?

For the first game:

Copy $ checkersd query checkers show-stored-game 1

This should print:
Copy afterIndex: "2" # The second game you created
beforeIndex: "-1" # No game
...

For the second game, run:

Copy $ checkersd query checkers show-stored-game 2

This should print:
Copy afterIndex: "-1" # No game
beforeIndex: "1" # The first game you created
...

Your FIFO in effect has the game IDs [1, 2].

Add a third game, and confirm that your FIFO is [1, 2, 3].
6

What happens if Alice plays a move in game 2, the game in the middle?

Copy $ checkersd tx checkers play-move 2 1 2 2 3 --from $alice
$ checkersd query checkers show-system-info

This should print:
Copy SystemInfo:
    fifoHeadIndex: "1" # The first game you created
    fifoTailIndex: "2" # The second game you created and on which Bob just played
    nextId: "4"
7

Is game 3 in the middle now?

Copy $ checkersd query checkers show-stored-game 3

This should print:
Copy storedGame:
    afterIndex: "2"
    beforeIndex: "1"
    ...

Your FIFO now has the game IDs [1, 3, 2]. You see that game 2, which was played on, has been sent to the tail of the FIFO.
8

What happens if Alice rejects game 3?

Copy $ checkersd tx checkers reject-game 3 --from $alice
$ checkersd query checkers show-system-info

This prints:
Copy SystemInfo:
    fifoHeadIndex: "1"
    fifoTailIndex: "2"
    nextId: "4"

There is no change because game 3 was in the middle, so it did not affect the head or the tail.

Fetch the two games by running the following two queries :

Copy $ checkersd query checkers show-stored-game 1

This prints:
Copy storedGame:
    afterIndex: "2"
    beforeIndex: "-1"
...

And:

Copy $ checkersd query checkers show-stored-game 2

This prints:
Copy storedGame:
    afterIndex: "-1"
    beforeIndex: "1"
...

Your FIFO now has the game IDs [1, 2]. Game 3 was correctly removed from the FIFO.
synopsis

To summarize, this section has explored:

    The use of a First-In-First-Out (FIFO) data structure to sort games from the least recently played at the top of the list to the most recently played at the bottom, in order to help identify inactive games which may become candidates for forced termination, which reduces undesirable and wasteful data stored on the blockchain.
    How forced termination of games is beneficial should you implement a wager system, as it prevents any assigned value from becoming locked into inactive games by causing the inactive player to forfeit the game and lose their wager.
    How any code solution which searches the entire data store for inactive games is computationally expensive, needlessly accessing many active games to identify any inactive minority (which may not even exist).
    How a FIFO data structure definitionally orders games such that inactive games rise to the top of the list, meaning code solutions can simply run until encountering the first active game and then stop, conserving gas fees.
    What new information and functions need to be added to your code; how to integrate them into the message handlers; how to update your unit tests to prevent them from failing due to these changes; and what tests to run to test the code.
    How to interact with the CLI to check the effectiveness of your new commands.



**Keep an Up to date Deadline**

- Make sure you have everything you need before proceeding:
    - You understand the concepts of Protobuf
    - Go is installed
    - You have the checkers blockchain codebase with the game FIFO. If not, follow the previous steps or check out the relevant version.

- In this section, you will:
    - Implement a deadline.
    - Work with dates.
    - Extend your unit tests.

- In the previous section you introduced FIFO that keeps the *oldest* games at it head and the most recently updated games at its tail.

- Just because a game has not been updated in a while does not mean that it has expired. To ascertain this you need to add a new field in the game.

**New Information**

- To prepare the field, add in the `StoredGame`'s protobuf definiton:

```
message StoredGame {
    ...
    string deadline = 9;
}
```

- To have Ignite CLI and Protobuf recompile this file, use:

`ignite generate proto-go`

- On each update the deadline will always be *now* plus a fixed duration. In this context, *now* refers to the block's time. Declare this duration as a new constant, plus how the date is to be represented - encoded in the saved game as a string:

```
const (
    MaxTurnDuration = time.Duration(24 * 3_600 * 1000_000_000) // 1 day
    DeadlineLayout = "2006-01-02 15:04:05:05.999999999 +0000 UTC"
)
```

**Date Manipulation** 

- Helper functions can encode and decode the deadline in the storage.

1. Define a new error:

```go
var (
    ...
    ErrInvalidDeadline = sdkerrors.Register(ModuleName, 1109, "deadline cannot be parsed: %s")
)
```

2. Add your date helpers. A reasonable location to pick is `full_game.go`:

```go
func (storedGame *StoredGame) GetDeadlineAsTime() (deadline time.Time, err error) {
    deadline, errDeadline := time.Parse(DeadlineLayout, storedGame.Deadline)
    return deadline, sdkerrors.Wrapf(errDeadline, ErrInvalidDeadline.Error(), storedGame.Deadline)
}

func FormatDeadline(deadline time.Time) string {
    return deadline.UTC().Format(DeadlineLayout)
}
```

- Note that `sdkerrors.Wrapf(err, ...)` conveniently returns `nil` if `err` is `nil`.

3. At the same time, add this to the `Validate` function:

```
...
_, err = storedGame.ParseGame()
if err != nil {
    return err
}
_, err = storedGame.GetDeadlineAsTime()
return err
```

**Updated deadline** 

- Next, you need to update this new field with its appropriate value:

1. At creation, in the message handler for game creation:

```
...
storedGame := types.StoredGame {
    ...
    Deadline: types.FormatDeadline(types.GetNextDeadline(ctx)),
}
```

2. After a move, in the message handler:

```
...
storedGame.MoveCount++
storedGame.Deadline = types.FormatDeadline(types.GetDeadline(ctx))
...
```

- Confirm that your project still compiles:

`ignite chain build`

**Unit Tests**

- After these changes, your previous unit tests fail. Fix them by adding `Deadline` whenever it should be. Do not forget that the time is taken from the block's timestamp. In the case of tests, it is stored in the context's `ctx.BlockTime()`. In effect, you need to add this single line: 

```
ctx := sdk.UnwrapSDKContext(context)
...
    require.EqualValues(t, types.StoredGame{
        ...
        Deadline: types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
    }, game)
```

Interact via the CLI

- There's not much to test here. Remember that you added a new field, but if your blockchain state already contains games then they are missing the new field:

`checkersd query checkers show-stored-game 1`

- This demonstrates some missing information:

```
...
deadline: ""
...
```

- In effect, your blockchain state is broken. Examine the section on migrations to see how to update your blockchain state to avoid such a breaking change. This broken state still lets you test the update of the deadline on play:

```
checkersd tx checkers play-move 1 1 2 2 3 --from $alice
checkersd query checkers show-stored-game 1
```

- This contains:

```
...
deadline: 2022-02-05 15:26:26.832533 +0000 UTC
...
```

- In the same vein, you can create a new game and confirm it contains the deadline.

**Synopsis**

- To summarize, this section has explored:
    - How to implement a new deadline field and work with dates to enable the application to check whether games which have not been recnetly updated have expired or not.
    - How the deadline must use the block's time as its reference point, since a non-deterministic `Date.now()` would change with each execution.
    - How to test your code to ensure that it functions as desired.
    - How to interact with the CLI to create a new game with the deadline field in place
    - How, if your blockchain contains preexisting games, that the blockchain state is now effectively broken, since the deadline field of those games demonstrates missing information (which can be corrected through migration).


**Record the Game Winner** 

- Make sure  you have everything you need before proceeding:
    - You understand the concepts of Protobuf
    - Go is installed
    - You have the checkers blockchain codebase with a deadline field and its handling. If not, follow the previous steps or check out the relevant version.


- In this section, you will:
    - Check for a game winner.
    - Extend unit tests.


- To be able to terminate games, you need to discern between games that are current and thoes that have reached an end - for example,  when they have been won.
- Therefore a good field to add is for the **winner**. It needs to contain:
    - The winner of a game that reaches completion.
    - Or winner by *forfeit* when a game has expired.
    - Or a neutral value when the game is active.

- In this exercise a draw is not handled and it would perhaps require yet another value to save in *winner*.

- It is time to introduce another consideration. When a gam has been won, no one else is going to play it. Its board will no longer be updated adn is no longer used for any further decisions. In effect, the board becomes redundant. With a view to keeping a node's storage requirement low, you should delete the board's content but keep the rest of the game's information.

- To keep a trace of the last state of the board, you emit it with an event.

**New Information** 

- In the `StoredGame` Protobuf definition file:

```go
message StoredGame {
    ...
    string winner = 10;
}
```

- Have ignite clie and protobuf recompile this file:

`ignite generate proto-go`


- Add a helper function to get the winner's address, if it exists. A good location is in `full_game.go`

```
func (storedGame StoredGame) GetPlayerAddress(color string) (address sdk.AccAddress, found bool, err error) {
    black, err := storedGame.GetBlackAddress()
    if err != nil {
        return nil, false, err
    }
    red, err := storedGame.GetRedAddress()
    if err != nil {
        return nil, false, err
    }
    // make a map and then search it pretty interesting.
    address, found = map[string]sdk.AccAddress{
        rules.PieceStrings[rules.BLACK_PLAYER]: black,
        rules.PieceStrings[rules.RED_PLAYER]: red,
    }[color]
    return address, found, nil
}

func (storedGame StoredGame) GetWinnerAddress() (address sdk.AccAddress, found bool, err error) {
    return storedGame.GetPlayerAddress(storedGame.Winner)
}
```

Maps.

- A Go map looks like this:

`map[KeyValue]ValueType`

- where `KeyType` may be any type that is comparable (more on this later), and ValueType may be any type at all, including another map!

- This variable `m` is a map of string keys to int values:

`var m map[string]int`


- Map types are reference types, like pointers or slices, and so the value of m above is nil; it doesn't point toa n initialized map. A nil map behaves like an empty map when reading, but attempts to write a nil map will cause a runtime panic; don't do that. To initialize a map, use the built in make function:

`m = make(map[string]int)`

- The **make** function allocates and initializes a hash map data structure and returns a map value that points to it. The specifics of taht adata structure are an implementation detail of hte runtime and are not specified by the language itself. In this article we will focus on the *use* of maps, not their implementation.


**Working with Maps** 

- Go provides a familiar syntax for working with maps. this statement sets the key "route" to the value 66:

`m["route"] = 66`

- This statement retrieves the value stored under the key "route" and assigns it to a new variable i:

`i := m["route"]`

- If the requested key doesn't exist, we get the value type's *zero value*. In this case the value type is int, so the zero value is 0:

`j := m["root"]`
// j == 0

- The built in `len` function returns on the number of items in a map:

`n := len(m)`

- The built in delete function removes an entry from the map:

`delete(m, "route")`

- The delete function doesn't return anything, and will do nothing if the specified key doesn't exist.
- A two-value assignment tests for the existence of a key:

`i, ok := m["route"]`

- In this statement, the first value (i) is assigned the values stoed under the key "route". If that key doesnt exist, i is the value type's zero value (0). The second value (ok) is a bool that is true if the key exists in the map, and false if not.

- To test for a key without retrieving the value, use an underscore in place of the first value:

`_, ok := m["route"]`

- To iterate over the contents of a map, use the range keyword:

```
for key, value := range m {
    fmt.Println("Key:", key, "Value:", value)
}
```

- To initialize a map with some data, use a map literal:

```
commits := map[string]int {
    "rsc": 3711,
    "r": 2138,
    "gri": 1908,
    "adg": 912,
}
```

- The same syntax may be used to initialize an empty map, which is fundamentally identical to using the make function:

`m = map[string]int{}`

**Exploiting Zero Values** 

- It can be convenient that a map retrieval yields a zero value when the key is not present.
- For instance, a map of boolean values can be used as a set-like data structure (recall that the zero value for the boolean type is false). This example traverses a linked list of Nodes and prints their values. It uses a map of Node pointers to detect cycles in the list.


```
type Node struct {
    Next *Node
    Value interface{}
}

var first *Node

visited := make(map[*Node]bool)

for n := first; n != nil; n = n.Next {
    if visited[n] {
        fmt.Println("cycle detected")
        break
    }
    visited[n] = true
    fmt.Println(n.Value)
}
```

- Create node, loop through nodes. Right, if the the first node is visited already then print the cycle is detected and break out of teh funciton.
- If the node has not been visited mark as visited and print the value. 

- The expression `visited[n]` is true if n has been visted, or false if n is not present. There's no need to use the two value form to test for the presence of the  n in the map; the zero value default does it for us.



-----

**Update and check for the winner**

- The is a two-part update. You set the winner where relevant, but you also introduce new checks so that a game with a winner cannot be acted upon.

- Start with a new error that you define as a constant:


```
var (
    ...
    ErrGameFinished = sdkerrors.Register(ModuleName, 1110, "game is already finshed")
)
```

- Ans a new event attribute:

```
const (
    MovePlayedEventType = "move-played"
    ...
    MovePlayedEventBoard = "board"
)
```


- At creation, in the *create game* message handler, start with a neutral value:

```
...
storedGame := types.StoredGame {
    ...
    Winner: rules.PieceStrings[rules.NO_PLAYER],
}
```

- With further checks when handling a play in the handler:

1. Check that the gam has not finished yet:

```
    ...
    if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
        return nil, types.ErrGameFinished
    }
    isBlack := storedGame.Black = msg.Creator
    ...
```

2. Update the winner field, which remains neutral if there is no winner yet:

```
    ...
    storedGame.Winner = rules.PieceStrings[game.Winner()]
    systemInfo, found := k.Keeper.GetSystemInfo(ctx)
    ...
```

3. Handle the FIFO differently depending on whether the game is finished or not, and adjust the board:

```
+  lastBoard := game.String()
+  if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
+      k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)
+      storedGame.Board = lastBoard
+  } else {
+      k.Keeper.RemoveFromFifo(ctx, &storedGame, &systemInfo)
+      storedGame.Board = ""
+  }

```



4. Add the new attribute in the event:

```
...
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(types.MovePlayedEventType,
            ...
            sdk.NewAttribute(types.MovedPlaydEventWinner, rules.PieceStrings[game.Winner()]),
            sdk.NewAttribute(types.MovedPlayedEventBoard, lastBoard),
            ),
    )
...
```


- And when rejecting a game, in its handler:

```
...
    if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
        return nil, types.ErrGameFinished
    }
    if storedGame.Black == msg.Creator {
        ...
    }

```


- Confirm the code compiles, add unit tests, and you are ready to handle the expiration of games.


**Unit Tests** 


- Add tests for your new functions.

- You also need to update your existing tests so that they pass with a new Winner value. Most of your tests you need to add this line.

```
    ...
    require.EqualValues(t, types.StoredGame{
        ...
+      Winner:    "*",
    }, game1)
    ...

```

- This `"*"` means that your tests no games have reached a conclusion with a winner. Time to fix that. In a dedicated `full_game_helpers.go` file, prepare all the moves that will be played in the test. For convenience, a move will be written as:


```
type GameMoveTest struct {
    player string
    fromX  uint64
    fromY  uint64
    toX    uint64
    toY    uint64
}
```

- If you do not want to create a complete game yourself, you can choose this one:

```
var (
    Game1Moves = []GameMoveTest{
        {"b", 1, 2, 2, 3}, // "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"
        {"r", 0, 5, 1, 4}, // "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"
        {"b", 2, 3, 0, 5}, // "*b*b*b*b|b*b*b*b*|***b*b*b|********|********|b*r*r*r*|*r*r*r*r|r*r*r*r*"
        ...
        {"r", 3, 6, 2, 5}, // "*b*b****|**b*b***|*****b**|********|********|**r*****|*B***b**|********"
        {"b", 1, 6, 3, 4}, // "*b*b****|**b*b***|*****b**|********|***B****|********|*****b**|********"
    }
)
```

- You may want to add a small function that converts `"b"` and `"r"` into their respective player addresses:

```
func GetPlayer(color string) string {
    if color == "b" {
        return Bob
    }
    return Carol
}
```


For Each loop (slice or array)

```
a := []string{"Foo", "Bar"}
for i, s := range a {
    fmt.Println(i,s)
}
```

```
0 Foo
1 Bar
```
func PlayAllMoves(t *testing.T, msgServer types.MsgServer, context context.Context, gameIndex string, moves []GameMoveTest) {
    for _, move := range Game1Moves {
        _, err := msgServer.PlayMove(context, &types.MsgPlayMove{
            Creator:   GetPlayer(move.player),
            GameIndex: gameIndex,
            FromX:     move.fromX,
            FromY:     move.fromY,
            ToX:       move.toX,
            ToY:       move.toY,
        })
        require.Nil(t, err)
    }
}

For each move in game moves play move with the move to and from x and y.


- Now, in a new file, create the test that plays all the moves, and checks at the end that the game has been saved with the right winner and that the FIFO is empty again:

```
func TestPlayMoveUpToWinner(t *testing.T) {
    msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
    ctx := sdk.UnwrapSDKContext(context)

    testutil.PlayAllMoves(t, msgServer, context, "1", testutil.Game1Moves)

    systemInfo, found := keeper.GetSystemInfo(ctx)
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId:        2,
        FifoHeadIndex: "-1",
        FifoTailIndex: "-1",
    }, systemInfo)

    game, found := keeper.GetStoredGame(ctx, "1")
    require.True(t, found)
    require.EqualValues(t, types.StoredGame{
        Index:       "1",
        Board:       "",
        Turn:        "b",
        Black:       bob,
        Red:         carol,
        MoveCount:   uint64(len(testutil.Game1Moves)),
        BeforeIndex: "-1",
        AfterIndex:  "-1",
        Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
        Winner:      "b",
    }, game)
    events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
    require.Len(t, events, 2)
    event := events[0]
    require.Equal(t, event.Type, "move-played")
    require.EqualValues(t, []sdk.Attribute{
        {Key: "creator", Value: bob},
        {Key: "game-index", Value: "1"},
        {Key: "captured-x", Value: "2"},
        {Key: "captured-y", Value: "5"},
        {Key: "winner", Value: "b"},
        {Key: "board", Value: "*b*b****|**b*b***|*****b**|********|***B****|********|*****b**|********"},
    }, event.Attributes[(len(testutil.Game1Moves)-1)*6:])
}
```

- When checking the attributes, it only cares about the last five.

- Feel free to create another game won by the red player.


**Interact via the CLI**

- If you have created games in an earlier version of the code, you are now in a broke state. You cannot even play the old games because they have `.Winner == ""` and this will be caught by the `if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER]` test. Start again:

`ignite chain serve --reset-once`


- Do not forget to export `alice` and `bob` again, as explained in an earlier section under "interact via the CLI"

- Confirm that thtere is no winner for a game when created.


Alice
cosmos1af5vp38k2a2m5z6v47lv3pmd5xrtpt7yumetwn
Bob
cosmos1p3rdmmw535m3j9q48sdka344ux9a974xjxn2gr


```
checkersd tx checkers create-game $alice $bob --from $alice
checkersd query checkers show-stored-game 1
```

This should show:

```
...
    winner: '*'
...
```

- And when a player plays:

```
checkersd tx checkers play-move 1 1 2 2 3 --from $alice
checkersd query checkers show-stored-game 1
```

- This should show:

```
...
    winner: '*'
...
```

- Testing with the CLI up to the point where the game is resolved with a rightful winner is btter covered by unit tests or with a nice GUI. You will be able to partially test this in the next section, via a forfeit.


**Synopsis**

- To summarize, this section has explored:

    - How to prepare for terminating games by defining a winner field that differentiates between the outright winner of a completed game, the winner by forfeit when a game is expired, or a game which is still active.
    
    - What new information and functions to add and where, including the winner field, helper functions to get any winner's address, a new error for games already finished, and checks for various application actions.
    - How to update your tests to check the functionality of your code.
    
    - How interacting via the CLI is partially impeded by any existing test games now being in a broken state due to the absence of a value in the winner field, with recommendations for next actions to take.


**Auto-Expiring Games** 


- Make sure you have everything you need before proceeding:
    - You understand the concepts of ABCI
    - Go is installed
    - You have the checkers blockchain codebase with the elements necessary for forfeit. If not, follow the previous steps or check out the relevant sections.

- In this section, you will:
    - Do begin block and end block operations.
    - Forfeit games automatically.
    - Do garbage collection.

- In the previous section you prepared the experation of games:
    - A First-In-First-Out (FIFO) that always has old games at its head and freshly updated games at its tail.
    - A deadline field to guide the expiration.
    - A winner field to further assist with forfeiting.


**New Information** 

- A game expires in two different situations:
    1. It was never really played, so it was removed quietly. That includes a single move by a single player.
    2. Moves were played by both players, making it a proper game, and forfeit is the outcome because a player then failed to play a move in time.

- In the latter case, you want to emit a new event which differentiates forfeiting a game from a win involving a move. Therefore you define new error constants:

```
const (
    GameForfeitedEventType      = "game-forfeited"
    GameForfeitedEventGameIndex = "game-index"
    GameForfeitedEventWinner    = "winner"
    GameForfeitedEventBoard     = "board"
)
```

**Putting callbacks in place** 

- When you use Ignite CLI to scaffold your module, it creates the `x/checkers/module.go` file with a lot of functions to accommodate your application. In particular, the function that **may** be called on your module on `EndBlock` is named `Endblock`:

```
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
    return []abci.ValidatorUpdate{}
}
```

- Ignite CLI left this empty. It is here that you add what you need done right before the block gets sealed. Create a new fil named `x/checkers/keeper/end_block_server_game.go` to encapsulate the knowledge about game expiry.

- Leave your function empty for now:

```
func (k Keeper) ForfeitExpiredGames(goCtx context.Context) {
    // TODO
}

```

in `x/checkers/module.go` update `EndBlock` with:

```go
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
    am.keeper.ForfeitExpiredGames(sdk.WrapSDKContext(ctx))
    return []abci.ValidatorUpdate{}
}
```

- This ensures that **if** your module's `EndBlock` function is called the expired games will be handled. For the **whole application to call your module** you hav to instruct it to do so. This takes place in `app/app.go`, where the applicaion is initialized with the proper order to call the `EndBlock` functions in different modules. In fact, yours has already ben placed at the end by Ignite:

```
app.mm.SetOrderEndBlockers(
    crisistypes.ModuleName,
    ...
    checkersmoduletypes.ModuleName,
)
```

- Your `ForfeitExpiredGames` function will now be called at the end of each block.

- Also prepare a new error:

```
var (
    ...
    ErrCannotFindWinnerByColor = sdkerrors.Register(ModuleName, 1111, "cannot find winner by color: %s"))
)
```

**Expire games handler** 

- With the callbacks in place, it is time to code the expiration properly.

**Prepare the main loop** 

- In `ForfeitExpiredGames`, it is a matter of looping through the FIFO, starting from the head, and handling games that are expired. You can stop at the first active game, as all thoes that come after are also active thanks to the careful updating of the FIFO.

**Expire games handler** 

- With the callbacks in place, it is time to code the expiration properly.

**Prepare the main loop** 

- In `ForfeitExpiredGames`, it is a matter of looping through the FIFO, starting from the head, and handling games that are expired. You can stop at the first active game, as all thoes that come after are also active thanks to the careful updating of the FIFO.
    1. Prepare useful information:
```
    ctx := sdk.UnwrapSDKContext(goCtx)

    opponents := map[string]string {
        rules.PieceStrings[rules.BLACK_PLAYER]: rules.PieceStrings[rules.RED_PLAYER]
        rules.PieceStrings[rules.RED_PLAYER]: rules.PieceStrings[rules.BLACK_PLAYER]
    }



```

 Expire games handler

With the callbacks in place, it is time to code the expiration properly.
#
Prepare the main loop

In ForfeitExpiredGames, it is a matter of looping through the FIFO, starting from the head, and handling games that are expired. You can stop at the first active game, as all those that come after are also active thanks to the careful updating of the FIFO.

    Prepare useful information:

Copy ctx := sdk.UnwrapSDKContext(goCtx)

opponents := map[string]string{
    rules.PieceStrings[rules.BLACK_PLAYER]: rules.PieceStrings[rules.RED_PLAYER],
    rules.PieceStrings[rules.RED_PLAYER]:   rules.PieceStrings[rules.BLACK_PLAYER],
}
x checkers keeper end_block_server_game.go
View source

Initialize the parameters before entering the loop:
Copy systemInfo, found := k.GetSystemInfo(ctx)
if !found {
    panic("SystemInfo not found")
}

gameIndex := systemInfo.FifoHeadIndex
var storedGame types.StoredGame
x checkers keeper end_block_server_game.go
View source

Enter the loop:
Copy for {
    // TODO
}
x checkers keeper end_block_server_game.go
View source

See below for what goes in this TODO.

After the loop has ended do not forget to save the latest FIFO state:

    Copy k.SetSystemInfo(ctx, systemInfo)
    x checkers keeper end_block_server_game.go
    View source

So what goes in the for { TODO }?
#
Identify an expired game

    Start with a loop breaking condition, if your cursor has reached the end of the FIFO:

Copy if gameIndex == types.NoFifoIndex {
    break
}
x checkers keeper end_block_server_game.go
View source

Fetch the expired game candidate and its deadline:
Copy storedGame, found = k.GetStoredGame(ctx, gameIndex)
if !found {
    panic("Fifo head game not found " + systemInfo.FifoHeadIndex)
}
deadline, err := storedGame.GetDeadlineAsTime()
if err != nil {
    panic(err)
}
x checkers keeper end_block_server_game.go
View source

Test for expiration:

    Copy if deadline.Before(ctx.BlockTime()) {
        // TODO
    } else {
        // All other games after are active anyway
        break
    }
    x checkers keeper end_block_server_game.go
    View source

Now, what goes into this if "expired" { TODO }?
#
Handle an expired game

    If the game has expired, remove it from the FIFO:

Copy k.RemoveFromFifo(ctx, &storedGame, &systemInfo)
x checkers keeper end_block_server_game.go
View source

Check whether the game is worth keeping. If it is, set the winner as the opponent of the player whose turn it is, remove the board, and save:
Copy lastBoard := storedGame.Board
if storedGame.MoveCount <= 1 {
    // No point in keeping a game that was never really played
    k.RemoveStoredGame(ctx, gameIndex)
} else {
    storedGame.Winner, found = opponents[storedGame.Turn]
    if !found {
        panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Turn))
    }
    storedGame.Board = ""
    k.SetStoredGame(ctx, storedGame)
}
x checkers keeper end_block_server_game.go
View source

Emit the relevant event:
Copy ctx.EventManager().EmitEvent(
    sdk.NewEvent(types.GameForfeitedEventType,
        sdk.NewAttribute(types.GameForfeitedEventGameIndex, gameIndex),
        sdk.NewAttribute(types.GameForfeitedEventWinner, storedGame.Winner),
        sdk.NewAttribute(types.GameForfeitedEventBoard, lastBoard),
    ),
)
x checkers keeper end_block_server_game.go
View source

Move along the FIFO for the next run of the loop:

    Copy gameIndex = systemInfo.FifoHeadIndex
    x checkers keeper end_block_server_game.go
    View source

For an explanation as to why this setup is resistant to an attack from an unbounded number of expired games, see the section on the game's FIFO.


**Unit Tests**

- How do you test something that is supposed to happen during the `EndBlock` event? You call the function that will be called within `EndBlock` (i.e. `Keeper.ForfeitExpiredGames`). Create a new test file `end_block_server_gam_test.go` for your tests. The situations that you can test are:

1. A game was never played, whil alone in the state or not. Or two games were never played. In this case, you need to confirm that th game was fully deleted, and that an event was emitted with no winners;

```go
func TestForfeitUnplayed(t *testing.T) {
    _, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
    ctx := sdk.UnwrapSDKContext(context)
    game1, found := keeper.GetStoredGame(ctx, "1")
    require.True(t, found)
    game1.Deadline = types.FormatDeadlin(ctx.BlockTime().Add(time.Duration(-1)))
    keeper.SetStoredGame(ctx, game1)
    keeper.ForfeitExpiredGames(context)

    _, found = keeper.GetStoredGame(ctx, "1")
    require.False(t, found)

    systemInfo, found := keeper.GetSystemInfo(ctx)
    require.True(t, found)
    require.EqualValues(t, types.SystemInfo{
        NextId: 2,
        FifoHeadIndex: "-1",
        FifoTailIndex: "-1",
    }, systemInfo)
    events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
    require.Len(t, events, 2)
    event := events[0]
    require.EqualValues(t, sdk.StringEvent{
        Type: "game-forfeited",
        Attributes: []sdk.Attribute{
            {Key: "game-index", Value: "1"},
            {Key: "winner", Value: "*"},
            {Key: "board",  Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
        },
    }, event)
}

```

When to use :=

- As others have explained, `:=` is for both decleration, assignment, and also for redecleration; and it guesses (*infers*) the variable's type automatically.

- **For example, foo := 32 is short-hand form of:**

```
var foo int
foo = 32

// OR: 
var foo int = 32

// OR:
var foo = 32

```

Rules for this 

1st Rule:
- You can't use := outside of `funcs`. It's because, outside a func, a statement should start with a keyword.

```
// no keywords below, illegal.
illegal := 42

// `var` keyword makes this statement legal.
var legal = 42

func foo() {
    alsoLegal := 42
    // reason: it's in a func scope.
}
```

2nd Rule:

- You can't use them twice (in the same scope):

```
legal := 42
legal := 42 // <-- error
```

- Because, := introduces "**a new variable**", hence using it twice does not redeclare a second variable, so it's illegal.


3rd Rule:

- You can use them for multi-variable declerations and assignments;

```
foo, bar := 42, 314
jazz, bazz := 22, 7
```

4th Rule (Redecleration):

- You can use them twice in "**multi-variable**" declerations, *if one of the variables is new*:

```
foo, bar    := someFunc()
foo, jazz   := someFunc() // <-- jazz is new
baz, foo    := someFunc() // <-- baz is new
```

- This is legal, because, you're not declaring all the variables, you're just reassigning new values to the existing variables at the same time. This is called *redecleration*.

5th Rule:

- You can use the short decleration to declare a variable in a newer scope even if that variable is already declared with the same name before

```
var foo int = 34

func some() {
    // because foo here is scoped to some func
    foo := 42 // <-- legal
    foo = 314 // <-- legal
}
```

- Here, `foo := 42` is legal, because, it declares `foo` in `some()` func's scope. `foo = 314` is legal, because, it just assigns a new value to `foo`.

6th rule:

- you can declare the smae name in short statement blocks like **if, for, switch**:

```
foo := 42
if foo := someFunc(); foo == 314 {
    // foo is scoped to 314 here
    // ...
}
// foo is still 42 here
```

- Because, `foo` in `if foo := ...`, only belongs to that `if` clause and it's in a different scope.

- **So, as a general rule**: if you want to easily declare a variable you can use :=, or, if you want to overwrite an existing variable, you can use `=`.


--- test section done- --

2. A games was played with only one move, while alone in the state or not. Or two games were played in this way. In this case, you need to confim that the game was fully deleted, and that an event was emitted with no winners.

3. A game was played with at last two moves, while alone in the sate or not. Or two games wrep layed in this way. In this case, you neede to confirm the game was not deleted, and instead that a winner was announced, including in events. 

- Note how all the attributes of an event of a given type (such as "game-forfeited") aggregate in a single array. The context is not reset on a new transaction, so when testing attributes you either have to compare the full array or take slices to compare what matters. 

**Interact via the CLI**

- Currently, the game expiry is one day in the future. This is too long to test with the CLI. Temporarily set it to 5 minutes:

`MaxTurnDuration = time.Duration(5 * 60 * 1000_000_000)) // 5 minutes`

- Avoid having games in the FIFO that expire in a day because of your earlier tests.

`ignite chain serve --reset-once`


Export your aliases again:

```
export alice=$(checkersd keys show alice -a)
export bob=$(checkersd keys show bob -a)
```

alice
cosmos1zdz7yhg076l7ex2n2zrq5wn3whs8kgxlvl3p9d

bob
cosmos1zsquwt8ks7c73xdvaayygehg5dkf98ptnd6ek9



- Create three games one minute apart. Have Alice play the middle one, and both Alice and Bob play the last one:


Create three games one minute apart. Have Alice play the middle one, and both Alice and Bob play the last one:
1

First game:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
2

Wait a minute, then create your second game and play it:

Copy $ checkersd tx checkers create-game $alice $bob --from $bob
$ checkersd tx checkers play-move 2 1 2 2 3 --from $alice
3

Wait another minute, then create your third game and play on it:

Copy $ checkersd tx checkers create-game $alice $bob --from $alice
$ checkersd tx checkers play-move 3 1 2 2 3 --from $alice
$ checkersd tx checkers play-move 3 0 5 1 4 --from $bob

Space each tx command from a given account by a couple of seconds so that they each go into a different block - by default checkersd is limited because it uses the account's transaction sequence number by fetching it from the current state.

- If you want to overcom this limitation, look at `checkersd`'s `--sequence` flag:

`checkersd tx checkers create-game --help`

- And at your account's current sequence. For instance:

`checkersd query account $alice --output json | jq -r '.sequence'`

- Which returns something like:

`9`


- With three games in, confirm that you see them all:

`checkersd query checkers list-stored-game`


- List them again after two, three, four, and five minutes. You should see games 1 and 2 disappear, and game 3 being forfeited by Alice, i.e. `red` bob wins


`checkersd query checkers show-stored-game 3 --output json | jq '.storedGame.winner'`

- This prints out "r"

- Confirm that the FIFO no longer references the removed games no the forfeited game:

- `checkersd query checkers show-system-info`

- This should show

```
SystemInfo:
    fifoHeadIndex: "-1"
    fifoTailIndex: "-1"
    nextId: "4"
```

**Synopsis** 

- To summarize, this section has explored:
    - How games can expire under two conditions:
        - When a game never really begins or only only one player makes an opening move, inw chich case iti s removed quietly; or when both player have particpated but one has since fialed to play a move in time, in which case the game is forfeited.
    - What new information and functions need to be created, and to update `EndBlock` to call the `ForfeitExpiredGames` function at the end of each block.
    - The correct coding for how to prepare the main loop through the FIFO, idenfiy an expired game, and handle an expired game.
    - How to test your code to ensure that it fundctions as desired.
    - How to interact with teh CLI to check the effectiveness of your code for handling expired games.



**Let Players Set a Wager** 


- In this section you will:
    - Add wager information (only).
    - Update unit tests.


- With the introduction of game expiry in the previous section and other features, you have now addressed the cases when two players start a game and finish it, or let it expire.

- In this section, you will go one step closer to adding an extra layer to a game, with wagers or stakes. Your application already includes all the necessary modules.

- Players choose to wager *money* or not, and the winner gets both wagers. The forfeiter loses their wager. To reduce complexity, start by letting players wager the staking token of your application.

- Now that no games can be left stranded, it is possible for players to safely wager on their games. How could this be implemented.

**Some initial thoughts** 

- When thinking about implementing a wager on games, ask:
    - What form will a wager take?
    - Who decides on the amount of wagers?
    - Where is a wager recorded?
    - At what junctures do you need to handle payments, refunds, and wins?

- This is a lot to go through. Therefore, the work is divided into two section. In this section, you only need to add new information, while the next section is where the tokens are actually handled.

- Some answers:
    - Even if only as a start, it makes sense to let the game creator decide on the wager.
    - It seems reasonable to save this information in the game itself so that wagers can be handled at any point in the lifecycle of the game.

**Code Needs**

- When it comes to your code: 
    - What Ignite CLI commands, if any, will assist you?
    - How do you adjust what Ignite CLI created for you?
    - Where do you make your changes?
    - What event should you emit?
    - How would you unit-test these new elements?
    - How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get?


**New Information** 

- Add this wager value to the `StoredGame`'s Protobuf definition:

```go
    message StoredGame {
        ...
        uint64 wager = 11;
    }
```
stored_game.proto

- You can let players choose the wager they want by adding a dedicated field in the message to create a game, in `proto/checkers/tx.proto`

- Have Ignite CLI and Protobuf recompile these two files:

`ignite generate proto-go`


- Now add a helper function to `StoredGame` using Cosmos SDK `Coin` in `full_game.go`:

```go
func (storedGame *StoredGame) GetWagerCoin() (wager sdk.Coin) {
    return sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(int64(storedGame.Wager)))
}
```
full_game.go


- This encapsulates information about the wager (where `sdk.DefaultBondDenom` is most likely `"stake"`)

**Saving the wager** 

- Time to ensure that the new field is saved in the storage and it is part of the creation event.

1. Define a new event key as a constant:

```
const (
    ...
    GameCreatedEventWager = "wager"
)
```
keys.go

2. Set the actual value in the new `StoredGame` as it is instantiated in the create game handler:

```
    storedGame := types.StoredGame{
        ...
        Wager: msg.Wager,
    }
```
msg_server_create_game.go


3. And in the event:

```
    ctx.EventManager().EmitEvent(
        sdk.NewEvent(sdk.EventTypeMessage,
        ...
        sdk.NewAttribute(types.GameCreatedEventWager, strconv.FormatUint(msg.Wager, 10)),
        )
    )

```

4. Modify the constructor among the interface definition of `MsgCreateGame` in `x/checkers/types/message_create_game.go` to avoid suprises:

```
func NewMsgCreateGame(creator string, red string, black string, wager uint64) *MsgCreateGame {
    return &MsgCreateGame {
        ...
        Wager: wager,
    }
}
``` 
message_create_game.go

5. Adjust the CLI client accordingly:


   func CmdCreateGame() *cobra.Command {
        cmd := &cobra.Command{
-          Use:   "create-game [black] [red]",
+          Use:   "create-game [black] [red] [wager]",
            Short: "Broadcast message createGame",
-          Args:  cobra.ExactArgs(2),
+          Args:  cobra.ExactArgs(3),
            RunE: func(cmd *cobra.Command, args []string) (err error) {
                argBlack := args[0]
                argRed := args[1]
+              argWager, err := strconv.ParseUint(args[2], 10, 64)
+              if err != nil {
+                  return err
+              }

                clientCtx, err := client.GetClientTxContext(cmd)
                if err != nil {
                    return err
                }
                msg := types.NewMsgCreateGame(
                    clientCtx.GetFromAddress().String(),
                    argBlack,
                    argRed,
+                  argWager,
                )
                if err := msg.ValidateBasic(); err != nil {
                    return err
                }
                return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
            },
        }
        flags.AddTxFlagsToCmd(cmd)
        return cmd
    }


**Interact via the CLI**

- With the tests done, see what happens at the command line. All there is to check at this stage is that the wager field appears where expected.

- After restarting the Ignite CLI, how much do Alice and Bob has to start with?


`checkersd query bank balances $alice`
`checkersd query bank balances $bob`

- This prints:

```
balances:
- amount: "100000000"
  denom: stake
- amount: "20000"
  denom: token
  pagination:
  next_key: null
  total: "0"
balances:
- amount: "100000000"
  denom: stake
- amount: "10000"
  denom: token
  pagination:
  next_key: null
  total: "0"
```

- Create a game with a wager:

`checkersd tx checkers create-game $alice $bob 1000000 --from $alice`

- Which mentions the wager:

```
...
raw_log: '[{"events":[{"type":"message","attributes":[{"key":"action","value":"create_game"}]},{"type":"new-game-created","attributes":[{"key":"creator","value":"cosmos1yysy889jzf4kgd84mf6649gt6024x6upzs6pde"},{"key":"game-index","value":"1"},{"key":"black","value":"cosmos1yysy889jzf4kgd84mf6649gt6024x6upzs6pde"},{"key":"red","value":"cosmos1ktgz57udyk4sprkpm5m6znuhsm904l0een8k6y"},{"key":"wager","value":"1000000"}]}]}]'
```

- Confirm that the balances of Alice and Bob are unchanged, as expected.

- Was the game stored correctly?

`checkersd query checkers show-stored-game 1`

- This returns:

```
storedGame:
    ...
    wager: "1000000"
```

- This confirms what you expected with regards to the command-line interactions.

**Synopsis**

- To summarize, this section has explored:
    - How to add the new "wager" valu, modify the "create a game" message to allow players to choose the wager they want to make, and add a helper function.
    - How to save the wager and adjust an event, modifiying the create game handler.
    - How to minimally adjust unit tests.
    - How to interact via the CLI to check that wager values are being recorded.


**Handle Wager Payments** 


- Make sure you have everything you need before proceeding:
    - You understand the concept of modules and keepers.
    - Go is installed.
    - you have the checkers blockchain codebase up to the game wager. If not, follow the previous steps or check out the relevant version.

- In this section, you will:
    - Work with the bank module.
    - Handle money.
    - Use mocks.
    - Add integration tests.


- In the previous section, you introduced a wager. On its own, having a `Wager` filed is just a piece of information, it does not transfer tokens just by existing.

- Transferring tokens is what this section is about.

**Some Initial Thoughts** 

- When thinking about implementing a wager on games, ask:
    - Is there any desirable atomicity of actions?
    - At what junctures do you need to handle payments, refunds, and wins?
    - Are there errors to report back?
    - What event should you emit?

- In the case of this example, you can consider that:
    - Although a game creator can decide on a wager, it should only be the holder of the tokens that can decide when they are being taken from their balance.
    - You might think of adding new message type, one that indicates that a player puts its wager in escrow. On the other hand, you can leverage the existing messages and consider that when a player makes their first move, this expresses a willingness to participate, and therfore the tokens can be transferred at this juncture.
    - For wins and losses, it is easy to imagine that hte code handles the payout at the time a game is resolved. 


**Code Needs**

- When it comes to your code:
    - What Ignite CLI commands, if any, will assist you?
    - How do you adjust what Ignite CLI created for you?
    - Where do you make your changes?
    - How would you unit-test these new elements?
    - Are unit tests sufficient here?
    - How would you use Ignite CLI to locally run a one-node blockchain and interact with it via the CLI to see what you get?

- Here are some elements of response:
    - Your module needs to call the bank to tell it to move tokens.
    - Your module needs to be allowed by the bank to keep tokens in escrow.
    - How would you test your module when it has such dependencies on the bank?


**What is to be done** 

- A lot is to be done. In order you will:
    -  Make it possible for your checkers module to call certain functions of the bank to transfer tokens.
    - Tell the bank to allow your checkers module to hold tokens in escrow.
    - Create helper functions that encapsulate some knowledge about when and how to transfer tokens.
    - Use the helper functions that encapsulate some knowledge about when and how to transfer tokens.
    - Use the helper functions at the right places in your code.
    - Update your unit tests and make use of mocks for that. You will create the mocks, create helper functions and use all that.
    - Prepare your code to accept integration tests.
    - Create helper functions that will make you integration tests more succinct.
    - Add integration tests that create a full app and test proper token bank balances.

**Declaring Expectations**

    - On its own the `Wager` field does not make the players pay the wager or recieve rewards. You need to add handling actions that ask the `bank` module to perform the required token transfers. For that, your keeper needs to ask for a `bank` instance during setup.

    - The only way to have access to a capability with the object-capability model of the Cosmos SDK is to be given the reference to an instance which already has this capability.

    - Payment handling is implemented by having your keeper hold wagers **in escrow** while the game is being played. the `bank` module has functions to transfer tokens from any account to your module and vice-versa.

    - Alternatively, your keeper could burn tokens instead of keeping them in escrow and mint them again when paying out. However, this makes your blockchain's total supply *falsely* fluctuate. Additionally, this burning and minting may prove questionable when you later introduce IBC tokens.

    - Declare an interface that narrowly declares the functions from other modules that you expect for your module. The conventional file for these declerations is `x/checkers/types/expected_keepers.go`.

    - The `bank` module has many capabilities, but all you need here are two functions. Your module already expects one function of the bank keeper: SpendableCoins. Instead of expanding this interface, you add a new one and *redeclare* the extra functions you need like so:

    ```go
    type BankEscrowKeeper interface {
        SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
        SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
    }
    ```
- These two functions must exactly match the functions declared in the `bank`'s keeper.go file. Copy the declerations directly from the `bank`'s file.
- In Go, any object with thes two functions is a `BankEscrowKeeper`.


**Obtaining the capability**

- With your requirements declared, it is time to make sure your keeper recieves a refecne to a bank keeper. First add a `BankEscrowKeeper` to your keeper in `x/checkers/keeper/keeper.go`:


```
type (
    Keeper struct {
        bank types.BankEscrowKeeper
        ...
    }
)
```
keeper.go

- This `BankEscrowKeeper` is your newly declared narrow interface. DO not forget to adjust the constructor accordingly:

```go
func NewKeeper (
    bank types.BankEscrowKeeper,
    ...
) *Keeper {
    return &Keeper {
        bank: bank,
        ...
    }
}
```

- Next, update where the constructor is called and pass a proper instance of `BankKeeper`. This happens in `app/app.go`:

```go
app.CheckersKeeper = *checkersmodulekeeper.NewKeeper(
    app.BankKeeper,
    ...
)
```
- This `app.BankKeeper` is a full `bank` keeper that also conforms to your `BankEscrowKeeper` interface.

- Finally, inform the app tha your checkers module is going to hold balances in escrow by adding it to the **whitelist** of permitted modules:

```
maccPerms = map[string][]string {
    ...
    checkersmoduletypes.ModuleName: nil,
}
```

- If you compare it to the other `maccperms` lines, the new line does not mention any `authtypes.Minter` or `authtypes.Burner`. Indeed `nil` is what you need to keep in escrow. For your information, the bank creates an *address* for your module's escrow account. When you have the full `app`, you can access it with:

```
import(
    "github.com/alice/checkers/x/checkers/types"
)
checkersModuleAddress := app.AccountKeeper.GetModuleAddress(types.ModuleName)
```

- On its own the Wager field does not make players pay or receive rewards. YOu need to add handling actions that ask the bank module to perform the required token transfers. For that, your keeper needs to ask for a `bank` instance during setup. So it can call the bank module to move the funds to the checkers module from escrow. 

object-capability model ...

**Preparing expected errors** 

- There are several new error situations that you can enumerate with new variables:


```
    var (
        ...
+      ErrBlackCannotPay    = sdkerrors.Register(ModuleName, 1112, "black cannot pay the wager")
+      ErrRedCannotPay      = sdkerrors.Register(ModuleName, 1113, "red cannot pay the wager")
+      ErrNothingToPay      = sdkerrors.Register(ModuleName, 1114, "there is nothing to pay, should not have been called")
+      ErrCannotRefundWager = sdkerrors.Register(ModuleName, 1115, "cannot refund wager to: %s")
+      ErrCannotPayWinnings = sdkerrors.Register(ModuleName, 1116, "cannot pay winnings to winner: %s")
+      ErrNotInRefundState  = sdkerrors.Register(ModuleName, 1117, "game is not in a state to refund, move count: %d")
    )
```

**Money handling steps**

- With the `bank` now in your keeper, it is time to have your keeper handle the money. Keep this concern in its own fil, as the functions are reused on play, reject, and forfeit.

- Create the new file, `x/checkers/keeper/wager_handler.go` and add three functions to collect a wager, refund a wager, and pay winnings.

```
func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error
func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame)
func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame)
```
- x/checkers/keeper/wager_handler.go

- The `Must` prefix in the function means that the transaction either takes place or a `panic` is issued. If a player cannot pay the wager, it is a usr-side error and the user must be informed of a failed transaction. If the module cannot pay, it means the escrow account has failed. This latter error is much more serious: an invariant may have been violated and the whole application must be terminated.

- Now set up collecting a wager, paying winnings, and refunding a wager:
    1. **Collecting wagers** happens on a player's first move. Therefore, differentiate between players:

```
    if storedGame.MoveCount == 0 {
        // Black plays first
    } else {
        // Red plays second
    }
    returns nil
```
x/checkers/keeper/wager_handler.go

- When there are no moves, get the address for the black player:


```
black, err := storedGame.GetBlackAddress()
if err != nil {
    panic(err.Error())
}
```
x/checkers/keeper/wager_handler.go


- Try to transfer into the escrow:

```
err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins)
```

- Do sam for rd player

2. **Paying Winnings** takes place when the game has declared winner.
- First get the winner. "No Winner" is **not** an acceptable situation in this `MustPayWinnings`. The caller of the function must ensure there is a winner:

```go
winnerAddress, found, err := storedGame.GetWinnerAddress()
if err != nil {
    panic(err.Error())
}
if !found {
    panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Winner))
}
```

- Calculate the winnings to pay:

```
winnings := storedGame.GetWagerCoin()
if storedGame.MoveCount == 0 {
    panic(types.ErrNothingToPay.Error())
} else if 1 < storedGame.MoveCount {
    winnings = winnings.Add(winnings)
}
```

- You double the wager only if the red player has also played and tehrefore both players have paid their wagers.

- If you did this wrongly, you could end up in a situation where a game with a single move pays out as if both players had played. This would be a serious bug that an attacker could exploit to drain your modules escrow fund.

- Then pay the winner: 

```
err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, winnerAddress, sdk.NewCoins(winnings))
if err != nil {
    panic(fmt.Sprintf(types.ErrCannotPayWinnings.Error(), err.Error()))
}
```

3. Finally, **refunding wagers** takes place when the game has partially started, i.e. only one party has paid, or when the game ends in a draw. In this narrow case of `MustRefundWager`:

```go
if storedGame.MoveCount == 1 {
    // Refund
} else if storedGame.MoveCount == 0 {
    // Do nothing
} else {
    // TODO Implement a draw mechanism.
    panic(fmt.Sprintf(types.ErrNotInRefundState.Error(), storedGame.MoveCount))
}
```

- Refund the black player when there has been a single move:

```go
black, err := storedGame.GetBlackAddress()
if err != nil {
    panic(err.Error())
}
err = k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, black, sdk.NewCoins(storedGame.GetWagerCoin()))
if err != nil {
    panic(fmt.Sprintf(types.ErrCannotRefundWager.Error(), err.Error()))
}
```

- You will notice that no special case is made when the wager is zero. This is a design choice here, and which way you choose to go is up to you. Not contacting the bank unnecessarily is cheaper in gas. On the other hand, why not outsource the zero check to the bank?


**Insert wager handling**

- With the desired steps defined in the wager handling functions, it is time to invoke them at the right places in the message handlers.
1. When a player plays for the first time:
    ```
    err = k.Keeper.CollectWager(ctx, &storedGame)
    if err != nil {
        return nil, err
    }
    ```
2. When a player wins as a result of a move:
    ```
    if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
        ...
    } else {
        ...
        k.Keeper.MustPayWinnings(ctx, &storedGame)
    }
    ```


When a player plays for the first time:
Copy err = k.Keeper.CollectWager(ctx, &storedGame)
if err != nil {
    return nil, err
}
x checkers keeper msg_server_play_move.go
View source

When a player wins as a result of a move:
Copy if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
    ...
} else {
    ...
    k.Keeper.MustPayWinnings(ctx, &storedGame)
}
x checkers keeper msg_server_play_move.go
View source

When a player rejects a game:
Copy k.Keeper.MustRefundWager(ctx, &storedGame)
x checkers keeper msg_server_reject_game.go
View source

When a game expires and there is a forfeit, make sure to only refund or pay full winnings when applicable. The logic needs to be adjusted:
Copy if deadline.Before(ctx.BlockTime()) {
    ...
    if storedGame.MoveCount <= 1 {
        ...
        if storedGame.MoveCount == 1 {
            k.MustRefundWager(ctx, &storedGame)
        }
    } else {
        ...
        k.MustPayWinnings(ctx, &storedGame)
        ...
    }
}
x checkers keeper end_block_server_game.go


**Unit tests**

- If you try running your existing tests you get a compilation error on the test keeper builder. Passing `nil` would not get you far with the tests and creating a fully-fledged bank keper would be a lot of work and not a unit test. See the integration tests below for that.

- Instead, you create mocks and use them in unit testss, not only to get the existing tests to pass but also to verify that the bank is called as expected.

**Prepare Mocks** 

- It is better to create some **mocks**. The Cosmos SDK does not offer mocks of its objects so you have to create your own. For that, the gomock library is a good resource. Install it:

- In a unit test, mock objects can simulate the behaviour of complex, real objects and are therefore useful when a real object is impractial or impossible to incorporate into a unit test.


`go install github.com/golang/mock/mockgen@v1.6.0`


Docker rebuild stuff

- create or rebuild docker image

`docker build -f Dockerfile-ubuntu . -t checkers_i`


https://tutorials.cosmos.network/hands-on-exercise/2-ignite-cli-adv/5-payment-winning.html#


https://tutorials.cosmos.network/hands-on-exercise/2-ignite-cli-adv/5-payment-winning.html#prepare-mocks


up to here...




























































Leverage Module
- This document specifies the x/leverage module of the Umee chain.
- The leverage module allows users to supply and borrow assets, and implements various features to support this, such as token accept-list, a dynamic interest rate module, incentivised liquidation and undercollateralized debt, and automatic reserve-based repayment of bad debt.

- The leverage module depends directly on x/oracle for asset prices, and interacts indirectly with x/uibc, x/gravity, adn the cosmos x/bank module.










- How would you collateralize NFTS? You have base nft levels like 10, 100, 1000 they would always be able to be transacted at these amounts. Like a base immutable value. but any more than that it is up to the market to decide.







Could escrow money in smart conract tehn have the backedd just call the contract which buys the ens so the money is not in some hot wallet or whatever, a fee is taken on the smart contract.
Only allows sniping of something within 30 days of the ens in the break in period or whatever.

Have like percentage of howm uch you want to pay in gas at the time like how urgent. 

will be automated and will work well, if it doesnt go through money will be returned without fee but yeah.








- The range expression, a, is **evaluated once** before beginning the loop.
- The iteration values are assigned to the respective iteration variables, i and s, **as in an assignment statement**.
- The second iteration variable is optional.
- For a nil slice, the number of iterations is 0.

- So the i is the iteration count, and the s is the value inside the array of strings. 






















**

NOTE: 
- Keeping the node's information low is a good thing so that there can be more validators and less storage on chain so deleting hte boards content would be good for space savings. 









































NEVER DOING IT FOR THE MONEY ALWAYS DOING IT FOR THE DREAM <3









Anything you have vision to do and are attracted to do you love, so if you cultivate love for the field you want to work with you will become connected to it and put your passion in it too. OYu love what you do .. How you cultivate obsession in a field you love it.  love is connection



How is this stored properly?? a question that is out of scope since now we are being more practical than theoretical

- Find out how to debug tests using docker ??

