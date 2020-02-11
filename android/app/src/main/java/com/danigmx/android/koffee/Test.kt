/**
 * Created by Daniel Gracia Machado on 2020-02-11.
 */
package com.danigmx.android.koffee

import android.util.Log
import okhttp3.*
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.RequestBody.Companion.toRequestBody
import org.json.JSONObject

class ReqLogin {
    public val email: String = "example@gmail.com"
    public val password: String = "123456example"

    fun json(): JSONObject {
        val json: JSONObject = JSONObject()
        json.put("email", email)
        json.put("password", password)
        return json
    }
}

class Test() {
    fun request(url: String) {
        var client = OkHttpClient()
        val a = ReqLogin().json()
        val mediaType = "application/json; charset=utf-8".toMediaType()
        val req: Request = Request.Builder()
            .url("${url}/api/user/login")
            .post(a.toString().toRequestBody(mediaType))
            .build();
        val response = client.newCall(req).execute()
        Log.d("tag", response.body.toString())
    }
}
