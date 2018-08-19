basic.showIcon(IconNames.Heart)
bluetooth.startTemperatureService()

bluetooth.startMagnetometerService()
bluetooth.onBluetoothConnected(() => {
    basic.showLeds(`
        . . . . .
        . # . # .
        . # # # .
        . # . # .
        . . . . .
        `)
    basic.pause(500)
    led.fadeOut()
})
