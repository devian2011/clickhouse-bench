## Clickhouse bench

```
==================== ClickHouse ========================
Operation: GroupByBrandIdLastDay
Times: 10
Detailed: 130.651584ms 77.839583ms 81.1585ms 103.267ms 94.79425ms 66.505333ms 79.015792ms 67.506ms 65.051125ms 67.270416ms
Avg Time: 83.305958ms
Operation: GroupByUserIdSumAmountLastThreeDays
Times: 10
Detailed: 131.882709ms 82.050916ms 127.150792ms 189.25475ms 151.670166ms 99.956375ms 125.806584ms 86.083917ms 81.414209ms 260.079292ms
Avg Time: 133.534971ms
Operation: SelectFourDays
Times: 10
Detailed: 3.7398045s 3.278264666s 4.0356475s 3.74048175s 3.495447667s 2.902541959s 3.307495833s 3.337894541s 2.889722208s 3.822022666s
Avg Time: 3.454932329s
=============================================================
```

## Postgres No Partition
```
==================== PostgresNoPartition ========================
Operation: GroupByBrandIdLastDay
Times: 10
Detailed: 620.729667ms 231.811083ms 265.152208ms 209.357333ms 263.614375ms 265.548083ms 302.189667ms 276.0315ms 154.467083ms 179.917459ms
Avg Time: 276.881845ms
Operation: GroupByUserIdSumAmountLastThreeDays
Times: 10
Detailed: 2.034595208s 1.35423125s 1.560466833s 1.223520167s 1.791847042s 1.428423584s 1.841467s 1.42721175s 1.659798167s 1.693377708s
Avg Time: 1.60149387s
Operation: SelectTwoDays
Times: 10
Detailed: 9.18794275s 9.030765166s 8.284673333s 7.369781541s 8.540127208s 8.51646325s 9.672986791s 7.526815584s 7.927654417s 9.043666958s
Avg Time: 8.510087699s
=============================================================
```

## Postgres With Partions

```
==================== PostgresPartition ========================
Operation: GroupByBrandIdLastDay
Times: 10
Detailed: 361.458875ms 64.525084ms 121.662667ms 211.155542ms 188.461583ms 104.092042ms 214.526167ms 108.020708ms 63.675584ms 41.67375ms
Avg Time: 147.9252ms
Operation: GroupByUserIdSumAmountLastThreeDays
Times: 10
Detailed: 609.726959ms 130.028375ms 383.881959ms 361.52825ms 311.560542ms 217.473084ms 353.216916ms 220.3925ms 146.19725ms 125.164916ms
Avg Time: 285.917075ms
Operation: SelectFourDays
Times: 10
Detailed: 6.054857417s 5.478374458s 5.3121955s 7.0487245s 5.359032917s 5.74572525s 5.859410041s 5.684535916s 5.495710875s 5.680156834s
Avg Time: 5.77187237s
=============================================================
```
