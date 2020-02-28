package com.danigmx.koffee

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import com.danigmx.koffee.koffee.R

class MainActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        val test = Test()
        test.request("http://localhost:8080")
    }

    // More to come...
    init {
        System.loadLibrary("keys")
    }

    external fun getMapsSdkApiKey()
}
