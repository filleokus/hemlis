# Hemlis

Hemlis is a CLI tool which generates an [age](https://github.com/FiloSottile/age) keypair and split the private key into shares using [Samir's secret sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing). They shares are encoded into words using the [Bytewords wordlist](https://github.com/BlockchainCommons/bc-bytewords) to make the shares easier to write or speak compared to a hex string.

![Hemlis](https://github.com/filleokus/hemlis/raw/main/example/screenshot.png "Hemlis")

## Features
- Built in, customizible, PDF generation
  - Generate with blank spaces and write the words manually with pen and paper or include them in the [PDF](https://github.com/filleokus/hemlis/blob/main/example/example-included-wordlist.pdf)
    - First method suggested unless you have access to a secure printer
  - Include / redact number of shares generated, threshold, and public key
    - If public key is included, a QR code of it is in the PDF
- "Checksum" of share included in PDF
  - Last 5 charachters of SHA256 hash included as identifier

## Installation
```
git clone https://github.com/filleokus/hemlis.git
cd hemlis
go build ./cmd/hemlis-combine && go build ./cmd/hemlis-gen
./hemlis-gen 
```

## Usage example
`hemlis-gen` generates the keypair and shares. `hemlis-combine` reconstructs it.

```
$ hemlis-gen -shares 3 -threshold 2
Public Key: age1t509ehdgvk27dpyq2wkevvj2m5nr67numtuzdkuavcx93c863cus6x3qjq
Private Key: AGE-SECRET-KEY-1LRARSQN92G28Z2L3E0944K6XNL7N2HEEJDDLEZGV3TS0GEGXRR4QUQA6CR
Number of shares: 3
Threshold: 2
--------------------------------
Share 1 (bcaf7)
|  0 part |  1 gyro |  2 gray |  3 iris |  4 vibe |
|  5 iced |  6 warm |  7 when |  8 luau |  9 roof |
| 10 ruby | 11 main | 12 very | 13 jump | 14 quiz |
| 15 slot | 16 jazz | 17 sets | 18 urge | 19 zaps |
| 20 when | 21 wolf | 22 gala | 23 fund | 24 yank |
| 25 away | 26 eyes | 27 tomb | 28 road | 29 jazz |
| 30 many | 31 dull | 32 legs |
--------------------------------
Share 2 (52f19)
|  0 down |  1 jowl |  2 buzz |  3 race |  4 able |
|  5 each |  6 dark |  7 iced |  8 kiln |  9 hill |
| 10 apex | 11 zaps | 12 drop | 13 huts | 14 holy |
| 15 owls | 16 rich | 17 real | 18 fact | 19 edge |
| 20 rich | 21 miss | 22 eyes | 23 iris | 24 brew |
| 25 keep | 26 ruby | 27 heat | 28 cash | 29 puma |
| 30 junk | 31 soap | 32 loud |
--------------------------------
Share 3 (d6312)
|  0 jazz |  1 beta |  2 gala |  3 frog |  4 epic |
|  5 dull |  6 into |  7 easy |  8 film |  9 navy |
| 10 judo | 11 exit | 12 solo | 13 free | 14 whiz |
| 15 lava | 16 lion | 17 quiz | 18 bias | 19 vial |
| 20 glow | 21 tomb | 22 next | 23 keep | 24 buzz |
| 25 webs | 26 pose | 27 liar | 28 fizz | 29 flew |
| 30 axis | 31 heat | 32 keno |
--------------------------------
Saving txt file
Generating PDF's
Print the PDF's and manually write the words on the papers
```
```
$ ./hemlis-combine -f shares.txt
✅ Share 0 (bcaf7) decoded
✅ Share 1 (d6312) decoded
Combined secret: AGE-SECRET-KEY-1LRARSQN92G28Z2L3E0944K6XNL7N2HEEJDDLEZGV3TS0GEGXRR4QUQA6CR
If not enough shares were used, the above secret will be incorrect

$ cat shares.txt
# First share
part
gyro
gray
iris
vibe
iced
warm
when
luau
roof
ruby
main
very
jump
quiz
slot
jazz
sets
urge
zaps
when
wolf
gala
fund
yank
away
eyes
tomb
road
jazz
many
dull
legs

# Third share
jazz
beta
gala
frog
epic
dull
into
easy
film
navy
judo
exit
solo
free
whiz
lava
lion
quiz
bias
vial
glow
tomb
next
keep
buzz
webs
pose
liar
fizz
flew
axis
heat
keno
```

# Dependencies
- `age` for key generation
  - Some vendored code for `bech32` encoding / decoding
- Samir's Secret Sharing algortihm from Hashicorp's Vault
- PDF generation from [johnfercher/maroto](https://github.com/johnfercher/maroto)
  - Known issue that the QR code does not render properly on (at least) Apple platforms: https://github.com/johnfercher/maroto/issues/413