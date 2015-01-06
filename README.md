nest
====
Command line interface to Nest.

### Build
```
go build -o /usr/local/bin/nest ./main.go
```

### Install
Download release and make it accessible in your path.

### Usage

#### Authorize

```
> nest

Login to nest using url: https://home.nest.com/login/oauth2?client_id=7348d6f1-1437-4935-a6e6-2a5567c96036&state=STATE

Enter pincode: #####

```

The access token will be printed and you can use this as environment variable for value NEST_ACCESS_TOKEN.

#### Show structures

```
> nest structures
Structures:
Id:             ##############
Name:           Home
Away:           home
```

#### Show thermostats

```
> nest thermostats
Thermostats:
Id:             ##############
Name:           Family Room
Name(long):     Family Room Thermostat
Temperature:    20.000000 C
Humidity:       40.000000
Status:         online
```

#### Get individual thermostat values

```
> nest --thermostat ### ambient-temperature
20
> nest --thermostat ### target-temperature
20
> nest --thermostat ### away-temperature-low
10
```

#### Set home / away

```
> nest -structure {structureid} home
Set to: home
```

### Shell alias 

```
export NEST_ACCESS_TOKEN=###
export THERMOSTATID=###
export STRUCTUREID=###

home() {
nest --structure "$STRUCTUREID" home
}

away() {
nest --structure "$STRUCTUREID" away
}

ambient-temperature() {
nest --thermostat "$THERMOSTATID" ambient-temperature
}

target-temperature() {
nest --thermostat "$THERMOSTATID" target-temperature
}

```

### Bonus

```
PS1='Current temperature: `nest --thermostat $THERMOSTATID ambient-temperature`:'

```

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2011-2014 Remco Verhoef.
Code released under [the MIT license](LICENSE).
