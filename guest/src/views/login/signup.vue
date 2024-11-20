<script setup lang="ts">
import { onMounted, reactive, watch } from "vue";
import { Hide, View, Message, Lock } from "@element-plus/icons-vue";
import Qq_color from "@/components/im_login/qq_color.vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import {SignupApi} from "@/api/user";

const router = useRouter();

const rules = reactive({
  email: [{ required: true, message: "Please input email", trigger: "blur" }],
  password: [
    { required: true, message: "Please input password", trigger: "blur" },
    { min: 8, max: 16, message: "Length should be 8 to 15", trigger: "blur" },
  ],
  confirmPassword: [
    {
      required: true,
      message: "Please input your password again",
      trigger: "blur",
    },
  ],
});

const form = reactive({
  email: "",
  password: "",
  confirmPassword: "",
  showPassword: false,
});

const togglePasswordVisibility = () => {
  form.showPassword = !form.showPassword;
};

const handleSignup = async () => {
  if (form.email && form.password && form.confirmPassword) {

    // 开始注册
    const res = await SignupApi(form)
    if(res.code){
      ElMessage.warning(res.msg)
      return
    }

    ElMessage.success("注册成功");
    await router.push("/login");
  } else {
    ElMessage.error("please input it to completion");
  }
};
</script>

<template>
  <div class="im_login">
    <div class="banner">
      <qq_color></qq_color>
    </div>
    <div class="login_form">
      <el-form :model="form" :rules="rules">
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            :prefix-icon="Message"
          >
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <!-- 使用 :type 动态绑定 input 类型 -->
          <el-input
            :type="form.showPassword ? 'text' : 'password'"
            v-model="form.password"
            placeholder="密码"
            :prefix-icon="Lock"
          >
            <!-- 添加一个 slot 来放置密码可见性切换按钮 -->
            <template #suffix>
              <i
                v-if="form.showPassword"
                style="cursor: pointer; font-size: 18px"
                @click="togglePasswordVisibility"
              >
                <el-icon>
                  <Hide />
                </el-icon>
              </i>
              <i
                v-else
                style="cursor: pointer; font-size: 18px"
                @click="togglePasswordVisibility"
              >
                <el-icon>
                  <View />
                </el-icon>
              </i>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <!-- 使用 :type 动态绑定 input 类型 -->
          <el-input
            :type="form.showPassword ? 'text' : 'password'"
            v-model="form.confirmPassword"
            placeholder="确认密码"
            :prefix-icon="Lock"
          >
            <!-- 添加一个 slot 来放置密码可见性切换按钮 -->
            <template #suffix>
              <i
                v-if="form.showPassword"
                style="cursor: pointer; font-size: 18px"
                @click="togglePasswordVisibility"
              >
                <el-icon>
                  <Hide />
                </el-icon>
              </i>
              <i
                v-else
                style="cursor: pointer; font-size: 18px"
                @click="togglePasswordVisibility"
              >
                <el-icon>
                  <View />
                </el-icon>
              </i>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button style="width: 100%" type="primary" @click="handleSignup"
            >注册
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style lang="scss">
.im_login {
  width: 500px;
  height: 406px;
  background-color: white;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 0 5px 3px rgba(0, 0, 0, 0.05);

  .banner {
    height: 140px;
    width: 100%;
    background-color: aquamarine;
  }

  .login_form {
    padding: 20px 60px;
  }
}
</style>
