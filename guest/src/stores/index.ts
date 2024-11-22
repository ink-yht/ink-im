import { defineStore } from "pinia";
import { parseToken } from "@/utils/parseToken";
import { type UserConfType, UserInfoApi } from "@/api/user";
import { ElMessage } from "element-plus";

interface userInfoType {
    exp: number;
    Uid: number;
    token: string;
    avatar: string;
}

const userInfo: userInfoType = {
    exp: 0,
    Uid: 0,
    token: "",
    avatar: "",
};

const userConfInfo: UserConfType = {
    id: 0,
    email: "",
    phone: "",
    nickname: "",
    abstract: "",
    avatar: "",
    friendOnline: false,
    sound: false,
    secureLink: false,
    savePwd: false,
    searchUser: 0,
    verification: 0,
}

export const useStore = defineStore("counter", {
    state: () => {
        return {
            userInfo: userInfo,
            userConfInfo: userConfInfo
        };
    },

    actions: {
        async setToken(token: string) {
            try {
                const payload = parseToken(token);
                this.userInfo.token = token;
                this.userInfo.exp = payload.exp;
                this.userInfo.Uid = payload.Uid;
                // 去拿用户的信息
                await this.getUserConf()
                this.userInfo.avatar = this.userConfInfo.avatar;
                // 持久化存储
                this.saveToken();
            } catch (e) {
                console.log(e);
            }
        },
        saveToken() {
            try {
                localStorage.setItem("UserInfo", JSON.stringify(this.userInfo));
            } catch (error) {
                console.error(error);
            }
        },
        loadToken() {
            const val = localStorage.getItem("UserInfo");
            if (!val) {
                // 没有登录，或者登录失效
                return;
            }
            try {
                this.userInfo = JSON.parse(val);
            } catch (e) {
                localStorage.removeItem("UserInfo");
                return;
            }
        },
        async getUserConf() {
            try {
                let res = await UserInfoApi();
                if (res.code) {
                    ElMessage.error(res.msg);
                    return;
                }
                this.userConfInfo = res.data
            } catch (e) {
                console.log(e)
            }
        }
    },

    getters: {
        isLogin(): boolean {
            try {
                return this.userInfo.token != "";
            } catch (e) {
                console.log(e);
            }
            return false
        },
    },
});
