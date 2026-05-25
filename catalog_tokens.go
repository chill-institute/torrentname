package torrentname

var (
	hdrCatalog = []aliasToken{
		{canonical: "HDR", aliases: []string{"HDR"}, patterns: []string{`HDR\b`}},
		{canonical: "HDR10+", aliases: []string{"HDR10+", "HDR10Plus", "HDR10P"}, patterns: []string{`HDR10Plus\b`, `HDR10P\b`, `HDR10\+`}},
		{canonical: "HDR10", aliases: []string{"HDR10"}, patterns: []string{`HDR10\b`}},
		{canonical: "DV", aliases: []string{"DV", "DoVi", "Dolby Vision"}, patterns: []string{`Dolby[ .-]+Vision`, `DoVi\b`, `DV\b`}},
		{canonical: "HLG", aliases: []string{"HLG"}, patterns: []string{`HLG\b`}},
	}

	resolutionCatalog = []aliasToken{
		{canonical: "480p", aliases: []string{"480P", "SD"}, patterns: []string{`480p`}},
		{canonical: "720p", aliases: []string{"720P", "HD"}, patterns: []string{`720p`}},
		{canonical: "1080p", aliases: []string{"1080P", "FHD"}, patterns: []string{`1080p`}},
		{canonical: "2160p", aliases: []string{"2160P", "4K", "UHD"}, patterns: []string{`2160p`, `4K`, `UHD`}},
	}

	qualityCatalog = []aliasToken{
		{canonical: "WEB-DL", aliases: []string{"WEB-DL", "WEB DL", "WEBDL", "PPV WEBDL", "HD WEBDL"}, patterns: []string{`(?:PPV|HD)[ .-]?WEB[ .-]?DL`, `WEB[ .-]?DL`}},
		{canonical: "WEBRip", aliases: []string{"WEBRip", "WEB Rip", "WBRip"}, patterns: []string{`WEB[ .-]?Rip`, `WBRip`}},
		{canonical: "WEB", aliases: []string{"WEB"}, patterns: []string{`WEB`}},
		{canonical: "BluRay", aliases: []string{"BluRay", "Blu Ray"}, patterns: []string{`Blu[ .-]?Ray`}},
		{canonical: "REMUX", aliases: []string{"REMUX", "BDRemux", "BluRay Remux", "Blu Ray Remux", "UHD BluRay Remux", "UHD Blu Ray Remux"}, patterns: []string{`UHD[ .-]+Blu[ .-]?Ray[ .-]+Remux`, `Blu[ .-]?Ray[ .-]+Remux`, `BDRemux`, `REMUX`}},
		{canonical: "HDRip", aliases: []string{"HDRip"}, patterns: []string{`HDRip`}},
		{canonical: "DVDRip", aliases: []string{"DVDRip", "DVDRIP"}, patterns: []string{`DVDRip`, `DVDRIP`}},
		{canonical: "BRRip", aliases: []string{"BRRip"}, patterns: []string{`BRRip`}},
		{canonical: "BDRip", aliases: []string{"BDRip"}, patterns: []string{`BDRip`}},
		{canonical: "HDTV", aliases: []string{"HDTV", "PPV.HDTV"}, patterns: []string{`(?:PPV\.)?HDTV`}},
		{canonical: "PDTV", aliases: []string{"PDTV", "PPV.PDTV"}, patterns: []string{`(?:PPV\.)?PDTV`}},
		{canonical: "TC", aliases: []string{"TC", "HDTC", "Telecine"}, patterns: []string{`(?:HD[ .-]?)?TC`, `Telecine`}},
		{canonical: "TS", aliases: []string{"HDTS", "Telesync"}, patterns: []string{`(?:HD[ .-]?)?TS`, `Telesync`}},
		{canonical: "CAM", aliases: []string{"HDCAM", "CAMRip"}, patterns: []string{`(?:HD)?CAM`, `CAMRip`}},
		{canonical: "DvDScr", aliases: []string{"DvDScr", "DVDScr"}, patterns: []string{`DvDScr`}},
	}

	codecCatalog = []aliasToken{
		{canonical: "x264", aliases: []string{"x264", "x 264", "x.264", "x-264", "x_264"}, patterns: []string{`x[ ._-]?264`}},
		{canonical: "x265", aliases: []string{"x265", "x 265", "x.265", "x-265", "x_265"}, patterns: []string{`x[ ._-]?265`}},
		{canonical: "H264", aliases: []string{"H264", "H 264", "H.264", "H-264", "H_264", "AVC"}, patterns: []string{`h[ ._-]?264`, `AVC`}},
		{canonical: "H265", aliases: []string{"H265", "H 265", "H.265", "H-265", "H_265", "HEVC"}, patterns: []string{`h[ ._-]?265`, `HEVC`}},
		{canonical: "AV1", aliases: []string{"AV1"}, patterns: []string{`AV1`}},
		{canonical: "XViD", aliases: []string{"XViD", "XVID"}, patterns: []string{`xvid`}},
	}

	languageCatalog = []aliasToken{
		{canonical: "MULTI", aliases: []string{"MULTI", "MultiLang"}, patterns: []string{`Multi(?:Lang)?`}},
		{canonical: "VOSTFR", aliases: []string{"VOSTFR"}, patterns: []string{`VOSTFR`}},
		{canonical: "VFF", aliases: []string{"VFF"}, patterns: []string{`VFF`}},
		{canonical: "VFQ", aliases: []string{"VFQ"}, patterns: []string{`VFQ`}},
		{canonical: "ENG", aliases: []string{"ENG", "English"}, patterns: []string{`ENG`, `ENGLISH`}},
		{canonical: "ITA", aliases: []string{"ITA", "Italian"}, patterns: []string{`ITA`, `ITALIAN`}},
		{canonical: "FRE", aliases: []string{"FRE", "French"}, patterns: []string{`FRE`, `FRENCH`}},
		{canonical: "GER", aliases: []string{"GER", "German"}, patterns: []string{`GER`, `GERMAN`}},
		{canonical: "SPA", aliases: []string{"SPA", "Spanish"}, patterns: []string{`SPA`, `SPANISH`}},
		{canonical: "LAT", aliases: []string{"LAT", "Latin"}, patterns: []string{`LAT`, `LATIN`}},
		{canonical: "RUS", aliases: []string{"RUS", "Russian"}, patterns: []string{`RUS`, `RUSSIAN`}},
		{canonical: "JPN", aliases: []string{"JPN", "JPS", "JAP", "Japanese"}, patterns: []string{`JPN`, `JPS`, `JAP`, `JAPANESE`}},
		{canonical: "UKR", aliases: []string{"UKR", "Ukrainian"}, patterns: []string{`UKR`, `UKRAINIAN`}},
		{canonical: "HIN", aliases: []string{"HIN", "Hindi"}, patterns: []string{`HIN`, `HINDI`}},
		{canonical: "KOR", aliases: []string{"KOR", "Korean"}, patterns: []string{`KOR`, `KOREAN`}},
		{canonical: "CHI", aliases: []string{"CHI", "Chinese"}, patterns: []string{`CHI`, `CHINESE`}},
	}

	editionCatalog = []aliasToken{
		{canonical: "Director's Cut", aliases: []string{"Director's Cut", "Directors Cut", "DC"}, patterns: []string{`Director'?s[ .-]+Cut`, `Directors[ .-]+Cut`, `DC`}},
		{canonical: "Hybrid", aliases: []string{"Hybrid"}, patterns: []string{`Hybrid`}},
		{canonical: "Theatrical Cut", aliases: []string{"Theatrical", "Theatrical Cut"}, patterns: []string{`Theatrical(?:[ .-]+Cut)?`}},
		{canonical: "Special Edition", aliases: []string{"Special Edition"}, patterns: []string{`Special[ .-]+Edition`}},
		{canonical: "Open Matte", aliases: []string{"Open Matte"}, patterns: []string{`Open[ .-]+Matte`}},
		{canonical: "Black and White", aliases: []string{"B&W", "BW", "Black and White"}, patterns: []string{`B&W`, `BW`, `Black[ .-]+and[ .-]+White`}},
		{canonical: "Dubbed", aliases: []string{"Dubbed"}, patterns: []string{`DUBBED`}},
		{canonical: "Dual Audio", aliases: []string{"Dual", "Dual Audio", "2Audio", "2Audios"}, patterns: []string{`DUAL`, `Dual[ .-]+Audio`, `2[ .-]*Audios?`}},
		{canonical: "Multi Subs", aliases: []string{"M Sub", "M Subs", "Multi Sub", "Multi Subs", "MultiSub"}, patterns: []string{`M[ .-]*Subs?`, `Multi[ .-]*Subs?`, `MultiSub`}},
	}

	audioCatalog = []aliasToken{
		{canonical: "MP3", patterns: []string{`MP3`}},
		{canonical: "TrueHD", patterns: []string{`TrueHD(?:[ .-]+(?:Atmos|` + audioChannelTokenPattern + `)){0,2}`}},
		{canonical: "DTS X", patterns: []string{`DTS[ .-]?X(?:[ .-]*(?:` + audioChannelTokenPattern + `))?`}},
		{canonical: "DTS-HD", patterns: []string{`DTS[ .-]?HD(?:[ .-]?(?:MA|HRA))?(?:[ .-]*(?:` + audioChannelTokenPattern + `))?`}},
		{canonical: "DTS", patterns: []string{`DTS(?:[ .-]*(?:` + audioChannelTokenPattern + `))?`}},
		{canonical: "EAC3", patterns: []string{`E-?AC-?3(?:[ .-]+Atmos)?(?:[ .-]*(?:` + audioChannelTokenPattern + `))?(?:[ .-]+Atmos)?`}},
		{canonical: "DDP Atmos", patterns: []string{`DDPA(?:[ .-]*(?:` + audioChannelTokenPattern + `))?`}},
		{canonical: "DDP", patterns: []string{`DDP(?:[ .-]*Atmos)?(?:[ .-]*(?:` + audioChannelTokenPattern + `))?(?:[ .-]+Atmos)?`}},
		{canonical: "DD+", aliases: []string{"DDPlus"}, patterns: []string{`DD(?:\+|Plus)(?:[ .-]+Atmos)?(?:[ .-]*(?:` + audioChannelTokenPattern + `))?(?:[ .-]+Atmos)?`}},
		{canonical: "DD", patterns: []string{`DD(?:[ .-]+Atmos)?(?:[ .-]*(?:` + audioChannelTokenPattern + `))?(?:[ .-]+Atmos)?`}},
		{canonical: "AAC", patterns: []string{`AAC[ .-]*LC`, `AAC[ .-]*(?:` + audioChannelTokenPattern + `)`, `AAC`}},
		{canonical: "AC3", patterns: []string{`AC3(?:[ .-]*(?:` + audioChannelTokenPattern + `)|\.5\.1)?`}},
		{canonical: "FLAC", patterns: []string{`FLAC(?:[ .-]*(?:` + audioChannelTokenPattern + `))?`}},
		{canonical: "PCM", patterns: []string{`L?PCM[ .-]*(?:` + audioChannelTokenPattern + `)`}},
		{canonical: "Opus", patterns: []string{`Opus[ .-]*(?:` + audioChannelTokenPattern + `)`}},
		{canonical: "Dual Audio", aliases: []string{"Dual-Audio"}, patterns: []string{`Dual[\- ]Audio`}},
		{canonical: "Channel", patterns: []string{`[257][ .][01]`, `[268][ .-]*CH`}},
		{canonical: "Atmos", patterns: []string{`Atmos`}},
		{canonical: "LiNE", patterns: []string{`LiNE`}},
	}

	containerCatalog = []aliasToken{
		{canonical: "MKV", aliases: []string{"mkv"}, patterns: []string{`MKV`}},
		{canonical: "AVI", aliases: []string{"avi"}, patterns: []string{`AVI`}},
		{canonical: "MP4", aliases: []string{"mp4"}, patterns: []string{`MP4`}},
	}

	extendedCatalog = []aliasToken{
		{canonical: "EXTENDED", aliases: []string{"EXTENDED CUT"}, patterns: []string{`EXTENDED(?:[ .-]?CUT)?`}},
	}

	hardcodedCatalog = []aliasToken{
		{canonical: "HC", aliases: []string{"HC SUB", "HC-SUB", "HCSUB", "KOR SUB", "KOR-SUB", "KORSUB"}, patterns: []string{`HC[ .-]?SUB`, `HCSUB`, `KOR[ .-]?SUB`, `KORSUB`, `HC`}},
	}

	properCatalog = []aliasToken{
		{canonical: "PROPER", patterns: []string{`PROPER`}},
	}

	repackCatalog = []aliasToken{
		{canonical: "REPACK", patterns: []string{`REPACK`}},
	}

	remasteredCatalog = []aliasToken{
		{canonical: "REMASTERED", patterns: []string{`REMASTERED`}},
	}

	widescreenCatalog = []aliasToken{
		{canonical: "WS", patterns: []string{`WS`}},
	}

	unratedCatalog = []aliasToken{
		{canonical: "UNRATED", patterns: []string{`UNRATED`}},
	}

	threeDCatalog = []aliasToken{
		{canonical: "3D", patterns: []string{`3D`}},
	}

	imaxCatalog = []aliasToken{
		{canonical: "IMAX", patterns: []string{`IMAX`}},
	}

	flagCatalog = mergeAliasTokens(extendedCatalog, hardcodedCatalog, properCatalog, repackCatalog, remasteredCatalog, widescreenCatalog, unratedCatalog, threeDCatalog, imaxCatalog)

	broadResolutionAliasContextCatalog = mergeAliasTokens(
		selectAliasTokens(qualityCatalog, "WEB-DL", "WEBRip", "WEB", "BluRay", "REMUX", "HDRip", "DVDRip", "BRRip", "BDRip", "HDTV", "PDTV", "DvDScr"),
		codecCatalog,
		hdrCatalog,
		audioCatalog,
	)
	broadResolutionReleaseContextPattern = tokenPatternAlternates(broadResolutionAliasContextCatalog)
	broadResolutionAliasContextPattern   = `(?:SD|HD|FHD)[ ._-]+(?:(?:` + broadResolutionReleaseContextPattern + `)|(?:` + sourceTokenPatternAlternates(sourceCatalog) + `)[ ._-]+(?:` + broadResolutionReleaseContextPattern + `))`

	sourceCatalog = []sourceToken{
		{canonical: "ABEMA", aliases: []string{"ABEMA"}, patterns: []string{`ABEMA`}},
		{canonical: "AMZN", aliases: []string{"AMZN"}, patterns: []string{`AMZN`}},
		{canonical: "ATVP", aliases: []string{"ATV", "ATV+", "ATVP"}, patterns: []string{`ATV\+?`, `ATVP`}},
		{canonical: "AUBC", aliases: []string{"AUBC"}, patterns: []string{`AUBC`}},
		{canonical: "BCORE", aliases: []string{"BCORE", "B CORE", "B-CORE", "B_CORE"}, patterns: []string{`B[ ._-]?CORE`}},
		{canonical: "BILI", aliases: []string{"BILI"}, patterns: []string{`BILI`}},
		{canonical: "CBC", aliases: []string{"CBC"}, patterns: []string{`CBC`}},
		{canonical: "CPNG", aliases: []string{"CPNG"}, patterns: []string{`CPNG`}},
		{canonical: "CR", aliases: []string{"CR"}, patterns: []string{`CR`}},
		{canonical: "CRAVE", aliases: []string{"CRAVE", "CRAV"}, patterns: []string{`CRAVE`, `CRAV`}},
		{canonical: "CRiT", aliases: []string{"CRiT", "CRIT"}, patterns: []string{`CRiT`}},
		{canonical: "DSNP", aliases: []string{"DSNP"}, patterns: []string{`DSNP`}},
		{canonical: "FOD", aliases: []string{"FOD"}, patterns: []string{`FOD`}},
		{canonical: "friDay", aliases: []string{"friDay", "FRIDAY"}, patterns: []string{`friDay`}},
		{canonical: "HAMI", aliases: []string{"Hami", "HAMI"}, patterns: []string{`Hami`}},
		{canonical: "HMAX", aliases: []string{"HMAX", "HBO MAX", "HBO-MAX", "HBO_MAX", "HBOM", "HBOMAX"}, patterns: []string{`HBO[ ._-]?MAX`, `HBOM`, `HMAX`}},
		{canonical: "HBO", aliases: []string{"HBO"}, patterns: []string{`HBO`}},
		{canonical: "HTSR", aliases: []string{"HTSR"}, patterns: []string{`HTSR`}},
		{canonical: "HULU", aliases: []string{"HULU"}, patterns: []string{`HULU`}},
		{canonical: "IQIY", aliases: []string{"iQIY", "IQIY"}, patterns: []string{`iQIY`}},
		{canonical: "ITVX", aliases: []string{"ITVX"}, patterns: []string{`ITVX`}},
		{canonical: "iT", aliases: []string{"iT", "iTUNES", "IT", "ITUNES"}, patterns: []string{`iTUNES`, `iT`}},
		{canonical: "KCW", aliases: []string{"KCW"}, patterns: []string{`KCW`}},
		{canonical: "KKTV", aliases: []string{"KKTV"}, patterns: []string{`KKTV`}},
		{canonical: "LINETV", aliases: []string{"LINE TV", "LINE-TV", "LINE_TV", "LINETV"}, patterns: []string{`LINE[ ._-]?TV`}},
		{canonical: "MAX", aliases: []string{"MAX"}, patterns: []string{`MAX`}, requiresContext: true},
		{canonical: "MY5", aliases: []string{"MY5"}, patterns: []string{`MY5`}},
		{canonical: "MYTVSUPER", aliases: []string{"MyTVSuper", "MYTVSUPER"}, patterns: []string{`MyTVSuper`}},
		{canonical: "NF", aliases: []string{"NF"}, patterns: []string{`NF`}},
		{canonical: "NOW", aliases: []string{"NOW"}, patterns: []string{`NOW`}, requiresContext: true, ambiguousGroup: true},
		{canonical: "OVID", aliases: []string{"OVID"}, patterns: []string{`OVID`}},
		{canonical: "PCOK", aliases: []string{"PCOK"}, patterns: []string{`PCOK`}},
		{canonical: "PLAY", aliases: []string{"PLAY"}, patterns: []string{`PLAY`}, requiresContext: true, ambiguousGroup: true},
		{canonical: "PMTP", aliases: []string{"PMTP"}, patterns: []string{`PMTP`}, ambiguousGroup: true},
		{canonical: "ROKU", aliases: []string{"ROKU"}, patterns: []string{`ROKU`}, requiresContext: true, ambiguousGroup: true},
		{canonical: "STAN", aliases: []string{"STAN"}, patterns: []string{`STAN`}, requiresContext: true, ambiguousGroup: true},
		{canonical: "STAR+", aliases: []string{"STAR", "STAR+"}, patterns: []string{`STAR\+?`}, requiresContext: true, ambiguousGroup: true},
		{canonical: "STRP", aliases: []string{"STRP"}, patterns: []string{`STRP`}, requiresContext: true, ambiguousGroup: true},
		{canonical: "TVING", aliases: []string{"TVING"}, patterns: []string{`TVING`}},
		{canonical: "TVER", aliases: []string{"TVER"}, patterns: []string{`TVER`}},
		{canonical: "UNEXT", aliases: []string{"U NEXT", "U-NEXT", "U_NEXT", "UNEXT"}, patterns: []string{`U[ ._-]?NEXT`}},
		{canonical: "VIKI", aliases: []string{"Viki", "VIKI"}, patterns: []string{`Viki`}},
		{canonical: "VIU", aliases: []string{"VIU"}, patterns: []string{`VIU`}},
		{canonical: "VRV", aliases: []string{"VRV"}, patterns: []string{`VRV`}},
		{canonical: "WAVVE", aliases: []string{"WAVVE"}, patterns: []string{`WAVVE`}},
		{canonical: "WETV", aliases: []string{"WETV"}, patterns: []string{`WETV`}},
		{canonical: "YOUKU", aliases: []string{"YOUKU"}, patterns: []string{`YOUKU`}},
		{canonical: "iP", aliases: []string{"iP", "IP"}, patterns: []string{`iP`}},
	}
)

var (
	hdrLookup        = buildAliasLookup(hdrCatalog)
	resolutionLookup = buildAliasLookup(resolutionCatalog)
	qualityLookup    = buildAliasLookup(qualityCatalog)
	codecLookup      = buildAliasLookup(codecCatalog)
	audioLookup      = buildAliasLookup(audioCatalog)
	containerLookup  = buildAliasLookup(containerCatalog)
	flagLookup       = buildAliasLookup(flagCatalog)
	languageLookup   = buildAliasLookup(languageCatalog)
	editionLookup    = buildAliasLookup(editionCatalog)
	sourceLookup     = buildSourceLookup(sourceCatalog)
)
