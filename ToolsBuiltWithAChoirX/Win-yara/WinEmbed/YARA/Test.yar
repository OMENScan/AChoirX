rule Connection_To_LocalHost {
    meta:
        author = "OmenScan"

    strings:
        $s1 = "127.0.0.1" ascii
        $s2 = "0.0.0.0" ascii

    condition:
        any of them
}
