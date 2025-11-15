# kwctl -- control a Kenwood TM-V71 (or similar) radio

This program allows you to control a Kenwood TM-V71 (and possibly similar models, like the TM-710) via the serial port.

## Commands

The following common options are available:

```
Usage of kwctl:
  -b, --bps int     bit rate (serial only) (default 57600)
  -d, --device string   serial device (default "/dev/radio0")
  -p, --pretty          pretty print output (default true)
  -v, --verbose count   increase logging verbosity
      --vfo string      select vfo on which to operate (default "1")
```

You can also use the following environment variables:

- `KWCTL_BPS` -- sets the default for the `--bps` option
- `KWCTL_DEVICE` -- sets the default for the `--device` option
- `KWCTL_PRETTY` set to `true` to enable pretty-print mode
- `KWCTL_VFO` -- sets the default for the `--vfo` option

### bands

```
Usage: kwctl bands [dual|single]

Get or set dual/single band mode.
```

### edit

```
Usage: kwctl edit [options] <channel>

Edit channel configuration.

Arguments:
        channel    Channel number to edit (0-999)

Options:
      --clear                 clear channel
      --copy int              copy data from another channel (default -1)
      --dcs dcs               DCS code (default 023)
      --lockout               skip channel during scan
      --mode mode             Mode (FM, NFM, AM) (default FM)
  -n, --name string           set channel name
      --no-lockout            don't skip channel during scan
      --no-reverse            disable reverse tx/rx
  -o, --offset frequencyMHz   offset in MHz (e.g., 0.6) (default 0.000000)
      --reverse               reverse tx/rx
  -r, --rxfreq frequencyMHz   frequency in MHz (e.g., 144.39) (default 0.000000)
      --rxstep stepSize       step size in hz (e.g., 5) (default 5)
      --rxtone tone           CTCSS tone when receiving (default 67.0)
  -s, --shift shift           Shift (simplex, up, down) (default simplex)
  -t, --tone-mode string      select tone mode (none, tone, tsql, dcs) (default "none")
      --txfreq frequencyMHz   frequency in MHz (e.g., 144.39) (default 0.000000)
      --txstep stepSize       step size in hz (e.g., 5) (default 5)
      --txtone tone           CTCSS tone when sending (default 67.0)
```

#### Examples

Configure a [repeater] on channel 90:

[repeater]: https://www.mmra.org/repeaters/BBY/index.html

```
$ kwctl -p edit 90 --rxfreq 146.820 --shift down --tone-mode tone --txtone 146.2 --offset 0.6 --name BAKBAY
┌────────┬────────┬────────────┬────────┬───────┬─────────┬──────┬───────┬───────┬──────────┬───────────┬─────────┬──────────┬──────┬──────────┬────────┬─────────┐
│ NAME   │ NUMBER │ RXFREQ     │ RXSTEP │ SHIFT │ REVERSE │ TONE │ CTCSS │ DCS   │ TONEFREQ │ CTCSSFREQ │ DCSCODE │ OFFSET   │ MODE │ TXFREQ   │ TXSTEP │ LOCKOUT │
├────────┼────────┼────────────┼────────┼───────┼─────────┼──────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┼──────────┼────────┼─────────┤
│ BAKBAY │ 090    │ 146.820000 │ 5      │ down  │ false   │ true │ false │ false │ 146.2    │ 67.0      │ 023     │ 0.600000 │ FM   │ 0.000000 │ 5      │ false   │
└────────┴────────┴────────────┴────────┴───────┴─────────┴──────┴───────┴───────┴──────────┴───────────┴─────────┴──────────┴──────┴──────────┴────────┴─────────┘
```

### channel

```
Usage: kwctl channel [options] [<channel>|up|down]

Get or set the current channel of the selected vfo.

Arguments:
        channel    Channel number (0-999) or 'up'/'down' to increment/decrement
```

### tune

```
Usage: kwctl tune [options]

Tune the selected VFO.

Options:
      --dcs dcs               DCS code (default 023)
  -f, --force                 change to vfo mode before tuning
      --mode mode             Mode (FM, NFM, AM) (default FM)
      --no-reverse            disable reverse tx/rx
  -o, --offset frequencyMHz   offset in MHz (e.g., 0.6) (default 0.000000)
      --reverse               reverse tx/rx
  -r, --rxfreq frequencyMHz   frequency in MHz (e.g., 144.39) (default 0.000000)
      --rxstep stepSize       step size in hz (e.g., 5) (default 5)
      --rxtone tone           CTCSS tone when receiving (default 67.0)
  -s, --shift shift           Shift (simplex, up, down) (default simplex)
  -t, --tone-mode string      select tone mode (none, tone, tsql, dcs) (default "none")
      --txtone tone           CTCSS tone when sending (default 67.0)
```

Note that if you can only tune the vfo when it is in vfo mode (which is why we have the `--force` option).

#### Examples

Show the current vfo configuration:

```
$ kwctl tune
1,145.090000,5,simplex,false,false,false,false,88.5,88.5,023,0.000000,FM
```

Or in pretty-print mode:

```
$ kwctl -p tune
┌─────┬────────────┬────────┬─────────┬─────────┬───────┬───────┬───────┬──────────┬───────────┬─────────┬──────────┬──────┐
│ VFO │ RXFREQ     │ RXSTEP │ SHIFT   │ REVERSE │ TONE  │ CTCSS │ DCS   │ TONEFREQ │ CTCSSFREQ │ DCSCODE │ OFFSET   │ MODE │
├─────┼────────────┼────────┼─────────┼─────────┼───────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┤
│ 1   │ 145.090000 │ 5      │ simplex │ false   │ false │ false │ false │ 88.5     │ 88.5      │ 023     │ 0.000000 │ FM   │
└─────┴────────────┴────────┴─────────┴─────────┴───────┴───────┴───────┴──────────┴───────────┴─────────┴──────────┴──────┘
```

Configure for use with a [repeater]

```
$ kwctl -p tune --rxfreq 146.820 --shift down --tone-mode tone --txtone 146.2 --offset 0.6
┌─────┬────────────┬────────┬───────┬─────────┬──────┬───────┬───────┬──────────┬───────────┬─────────┬──────────┬──────┐
│ VFO │ RXFREQ     │ RXSTEP │ SHIFT │ REVERSE │ TONE │ CTCSS │ DCS   │ TONEFREQ │ CTCSSFREQ │ DCSCODE │ OFFSET   │ MODE │
├─────┼────────────┼────────┼───────┼─────────┼──────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┤
│ 1   │ 146.820000 │ 5      │ down  │ false   │ true │ false │ false │ 146.2    │ 146.2     │ 023     │ 0.600000 │ FM   │
└─────┴────────────┴────────┴───────┴─────────┴──────┴───────┴───────┴──────────┴───────────┴─────────┴──────────┴──────┘
```

### version

Show version information.

### vfo

```
Usage: kwctl vfo [0|1]

Get or set ptt/control VFO.
```

Note that while the radio permits the PTT and control VFOs to be separate, this command always sets both to the same value.


### list

```
Usage: kwctl channel-list [options] <range> [<range> [...]]

List a range of channels.

Arguments:
        range      A range specification (e.g. "1", "1-10", "1,5,10,15,20")
```

#### Examples

```
$ kwctl list 1-4 10 11
[MRABBY] 001,146.820000,5,down,false,true,false,false,146.2,146.2,023,0.600000,FM,0.000000,5,false
[MRMDN ] 002,146.610000,5,down,false,true,false,false,146.2,146.2,023,0.600000,FM,0.000000,5,false
[MRAQCY] 003,146.670000,5,down,false,true,false,false,146.2,146.2,023,0.600000,FM,0.000000,5,false
[MRANRD] 004,146.715000,5,down,false,true,false,false,146.2,146.2,023,0.600000,FM,0.000000,5,false
[MRANRD] 010,446.775000,12.5,down,false,true,false,false,88.5,88.5,023,5.000000,FM,0.000000,5,false
[MRAHOP] 011,447.775000,12.5,down,false,true,false,false,88.5,88.5,023,5.000000,FM,0.000000,5,false
```

Or in pretty-print mode:

```
$ kwctl -p list 1-4
┌────────┬────────┬────────────┬────────┬───────┬─────────┬──────┬───────┬───────┬──────────┬───────────┬─────────┬──────────┬──────┬──────────┬────────┬─────────┐
│ NAME   │ NUMBER │ RXFREQ     │ RXSTEP │ SHIFT │ REVERSE │ TONE │ CTCSS │ DCS   │ TONEFREQ │ CTCSSFREQ │ DCSCODE │ OFFSET   │ MODE │ TXFREQ   │ TXSTEP │ LOCKOUT │
├────────┼────────┼────────────┼────────┼───────┼─────────┼──────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┼──────────┼────────┼─────────┤
│ MRABBY │ 001    │ 146.820000 │ 5      │ down  │ false   │ true │ false │ false │ 146.2    │ 146.2     │ 023     │ 0.600000 │ FM   │ 0.000000 │ 5      │ false   │
├────────┼────────┼────────────┼────────┼───────┼─────────┼──────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┼──────────┼────────┼─────────┤
│ MRMDN  │ 002    │ 146.610000 │ 5      │ down  │ false   │ true │ false │ false │ 146.2    │ 146.2     │ 023     │ 0.600000 │ FM   │ 0.000000 │ 5      │ false   │
├────────┼────────┼────────────┼────────┼───────┼─────────┼──────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┼──────────┼────────┼─────────┤
│ MRAQCY │ 003    │ 146.670000 │ 5      │ down  │ false   │ true │ false │ false │ 146.2    │ 146.2     │ 023     │ 0.600000 │ FM   │ 0.000000 │ 5      │ false   │
├────────┼────────┼────────────┼────────┼───────┼─────────┼──────┼───────┼───────┼──────────┼───────────┼─────────┼──────────┼──────┼──────────┼────────┼─────────┤
│ MRANRD │ 004    │ 146.715000 │ 5      │ down  │ false   │ true │ false │ false │ 146.2    │ 146.2     │ 023     │ 0.600000 │ FM   │ 0.000000 │ 5      │ false   │
└────────┴────────┴────────────┴────────┴───────┴─────────┴──────┴───────┴───────┴──────────┴───────────┴─────────┴──────────┴──────┴──────────┴────────┴─────────┘
```

### id

```
Usage: kwctl id

Display the radio ID response.
```

### mode

```
Usage: kwctl mode [vfo|memory|call|wx]

Get or set the operating mode for the selected VFO.
```

### txpower

```
Usage: kwctl txpower [high|medium|low]

Get or set the transmit power for the selected VFO.
```

## License

kwctl -- control a Kenwood TM-V71 (or similar) radio  
Copyright (C) 2025 Lars Kellogg-Stedman

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.

## Thanks

Thanks to LA3QMA for putting together https://github.com/LA3QMA/TM-V71_TM-D710-Kenwood, which was my source of information regarding available CAT commands.
