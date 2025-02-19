import { defineStore } from "pinia"
import { reactive, watch } from "vue"

type User = {
  autoSeen: boolean
  token: string
}

export const useUserStore = defineStore("user", () => {
  const user = reactive({
    autoSeen: false,
    token: "",
  })

  const localStorageName = "user"
  const userInStorage = localStorage.getItem(localStorageName)
  if (userInStorage) {
    const { autoSeen, token } = JSON.parse(userInStorage)
    user.autoSeen = autoSeen
    user.token = token
  }

  watch(
    () => user,
    (state) => {
      localStorage.setItem(localStorageName, JSON.stringify(state))
    },
    { deep: true }
  )

  const Login = (c: User) => {
    user.autoSeen = c.autoSeen
    user.token = c.token
  }

  const isLogin = () => {
    return !!user.token
  }
  const Logout = () => {
    (user.token = ""), (user.autoSeen = false)
  }

  return {
    user,
    Login,
    isLogin,
    Logout,
  }
})
