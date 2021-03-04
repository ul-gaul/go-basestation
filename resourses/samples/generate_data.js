'use strict';

setTimeout(() => {
    const START = performance.now()
    const count = 1000
    const rng = {
        alt: [rdm.gauss(2000, 1000), rdm.gauss(10000, 1000)],
        pressure: [25000, 80000],
    }
    for (let i in rng) {
        Object.defineProperties(rng[i], {
            min: { get() { return this[0] } },
            max: { get() { return this[1] } }
        })
    }

    const diff = {
        alt: (rng.alt.max - rng.alt.min) / (count / 2),
        pressure: (rng.pressure.max - rng.pressure.min) / (count / 2),
    }

    const step = (d) => rdm.gauss(d/2, d/4)
    const calc = (i, key) => i <= count/2
        ? rng[key].min + diff[key] * i + step(diff[key])
        : rng[key].max - diff[key] * (i-1-count/2) - step(diff[key])

    let data = []
    let ts = 0
    for (let i = 1; i <= count; i++) {
        ts += rdm.gauss(100, 50)
        data.push({
            time_stamp: Math.floor(ts),
            latitude: 0, // TODO
            longitude: 0, // TODO
            ns_indicator: rdm.bool('N', 'S'),
            ew_indicator: rdm.bool('E', 'W'),
            gps_satellites: rdm({ min: 15, max: 24, round: 0 }),
            altitude: calc(i, 'alt'),
            pressure: calc(i, 'pressure'),
            temperature: rdm.gauss(42),
            acceleration: {
                x: 0, // TODO
                y: 0, // TODO
                z: 0 // TODO
            },
            magnetometer: {
                x: 0, // TODO
                y: 0, // TODO
                z: 0 // TODO
            },
            angular_speed: {
                x: rdm.sinGauss({ idx: i, mod: .01, amplitude: 160, deviation: 5 }),
                y: rdm.sinGauss({ idx: i, shift: 2/3 * Math.PI/.1, mod: .01, amplitude: 175, deviation: 5 }),
                z: rdm.sinGauss({ idx: i, shift: 4/3 * Math.PI/.1, mod: .01, amplitude: 175, deviation: 5 }),
            },
            acquisition_board_state: {
                1: rdm.bool(1, 0),
                2: rdm.bool(1, 0),
                3: rdm.bool(1, 0)
            },
            power_supply_state: {
                1: rdm.bool(1, 0),
                2: rdm.bool(1, 0)
            },
            payload_board_state: {
                1: rdm.bool(1, 0)
            }
        })
    }

    for (let d of data) {
        for (let [k, v] of Object.entries(d)) {
            if (v && typeof v === 'object') {
                for(let [key, val] of Object.entries(v)) {
                    d[k + '_' + key] = val
                }
                delete d[k]
            }
        }
    }

    const headers = ('time_stamp,latitude,longitude,ns_indicator,ew_indicator,gps_satellites,altitude,pressure,temperature,'
        + ['acceleration','magnetometer','angular_speed'].map(k => [...'xyz'].map(c => [k,c].join('_')).join()).join() + ','
        + [1,2,3].map(k => `acquisition_board_state_${k}`).join() + ','
        + [1,2].map(k => `power_supply_state_${k}`).join() + ','
        + [1].map(k => `payload_board_state_${k}`).join()).split(',')

    Object.defineProperty(data, 'csv', {
        get() {
            let csv = headers.join(',')
            for (let row of this) {
                csv += '\n' + headers.map(k => row[k]).join()
            }
            return csv
        }
    })

    for (let header of headers) {
        if (header !== 'time_stamp') {
            Object.defineProperty(data, header, {
                get() {
                    return this.map((o) => ({ time_stamp: o.time_stamp, [header]: o[header] }))
                }
            })
        }
    }

    Object.defineProperty(data, 'show', {
        value(...keys) {
            return this.map((o) => {
                let ret = { time_stamp: o.time_stamp }
                for (let key of keys)
                    ret[key] = o[key]
                return ret
            })
        }
    })

    console.log(window.data = data)
})

function rdm({ min = 0, max = 1, round = 10 } = {}) {
    let v = (Math.random() * (max - min)) + min
    return round == 0 ? Math.round(v) : Math.round(v * (10 ** round)) / (10 ** round)
}

rdm.gauss = (m = 0, sd = 1) => {
    let y1, x1, x2, w;
    do {
        x1 = Math.random() * 2 - 1
        x2 = Math.random() * 2 - 1
        w = x1**2 + x2**2
    } while(w >= 1);
    w = Math.sqrt(-2 * Math.log(w) / 2);
    y1 = x1 * w;
    //y2 = x2 * w;
    return y1 * sd + m;
}

rdm.sinGauss = ({ idx, shift = 0, amplitude = 1, mod = 1, mean = 0, deviation = 1 }) =>
    rdm.gauss(amplitude * Math.sin(mod*idx + shift) + mean, deviation)

rdm.bool = (t = true, f = false) => Math.random() > 0.5 ? t : f;