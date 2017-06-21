# Alias

* `act`: Action
* `met`: Metadata
* `dat`: Data
* `fil`: File
    * `key`: File Key
    * `bin`: File Binary
    * `inf`: File Info
        * `siz`: File Size
        * `tolSiz`: Total File Size
        * `par`: Current Part
        * `tol`: Total Parts
* `suc`: Success Data
* `cod`: Code
* `err`: Error Data

# Dat

**`POST`** `ws://xpc.example.com/v1/notif`
**`POST`** `ws://xpc.example.com/v1/file`
**`POST`** `ws://xpc.example.com/v1/action`

```js
{
    act: "UploadVideo"
    met: {
        session: "AJjMC39xO1cpELfbGC8H4Os21G"
    },
    dat: {
        username: "YamiOdymel",
        title: "My Video!"
    },
    fil: {
        key: "AJjMC39xO1cpELfbGC8H4Os21G",
        inf: {
            siz: 937135,
            tolSiz: 3192351,
            par: 2,
            tol: 3
        },
        bin: "ÿØÿàJFIFÿÛC  %# , #&')*)-0-(0%()(ÿÛC   (((((((((((((((((((((((((((((((((((((((((((((((((((ÿÀTG"ÿÄ ÿÄµ}!1AQa"q2¡#B±ÁRÑð$3br %&'()*456789:CDEFGHIJSTUVWXYZcdefghijstuvwxyz¢£¤¥¦§¨©ª²³´µ¶·¸¹ºÂÃÄÅÆÇÈÉÊÒÓÔÕÖ×ØÙÚáâãäåæçèéêñòóôõö÷øùúÿÄ	ÿÄµw!1AQaq"2B¡±Á	#3RðbrÑ $4á%ñ&'()*56789:CDEFGHIJSTUVWXYZcdefghijstuvwxyz¢£¤¥¦§¨©ª²³´µ¶·¸¹ºÂÃÄÅÆÇÈÉÊÒÓÔÕÖ×ØÙÚâãäåæçèéêòóôõö÷øùúÿÚ?¥rÒ%Éèjj{À>Ñ.:f«úSÃÓ	ÍÒS)
iA¦H
-!©ÕXÉÍXZd²e<Ô«Ö SÍJ¦¥ ,#b¤V«©§RÐ2Ê
@¦SÍ"KQ·éT¡ëW#íQ!µIT*iÛ¸¬ÚòÔÆj2i ò(°Õ½*Â*¤dVcè*YH
<°ªqADRNäïÛæÅt»c#Þ¹K¶Ëí ´6ëV"5Pº¢\SR©ªÈÕ2Í¤Lô§½?<T7g­(¨³ÏzNiÀô§¨ö¦JÄP:Ó7Q¸R°;Ó"
    }
}
```

# Resp

```js
{
    met: {
    },
    dat: {
        key: "AJjMC39xO1cpELfbGC8H4Os21G"
    },
    // suc: {
    //  cod: "0000x0b",
    //  dat: {}
    // }
    err: {
      cod: "0000x0b",
      dat: {}
    }
}
```
