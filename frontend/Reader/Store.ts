import { defineStore } from "pinia"
import { ref } from "vue"

export type Category = {
  id: number
  name: string
  count: number
}

export const useCategoryStore = defineStore("category", () => {
  const list = ref([] as Category[])
  const selected = ref(0)

  const setCategories = (payload: { list: Category[]; index: number }) => {
    list.value = payload.list
    selected.value = payload.index
  }
  const setIndex = (index: number) => {
    selected.value = index
  }

  const nextCategory = () => {
    if (list.value.length == 0) {
      return
    }
    const index = selected.value + 1
    if (index === list.value.length) {
      selected.value = 0
    } else {
      selected.value = index
    }
  }
  const previousCategory = () => {
    if (list.value.length == 0) {
      return
    }
    const index = selected.value - 1
    if (0 > index) {
      selected.value = list.value.length - 1
    } else {
      selected.value = index
    }
  }

  const currentPointer = () => {
    const sel = list.value[selected.value]
    return sel ? sel.id : null
  }

  return {
    list,
    selected,
    setCategories,
    setIndex,
    nextCategory,
    previousCategory,
    currentPointer,
  }
})

export type ReadFlag = "Unseen" | "Seen" | "Setpin"

export type Entry = {
  date_epoch: number
  date: string
  description: string
  feed_id: number
  readflag: ReadFlag
  serial: number
  site_title: string
  subscription_id: number
  title: string
  url: string
}

type ContentSerialData = {
  serial: number
  feed_id: number
  readflag: ReadFlag
}

export const useContentStore = defineStore("content", () => {
  const list = ref([] as Entry[])
  const selected = ref(0)

  const setContents = (payload: { list: Entry[]; index: number }) => {
    list.value = payload.list
    selected.value = payload.index
  }
  const setIndex = (index: number) => {
    selected.value = index
  }
  const setReadflag = (flag: ReadFlag) => {
    list.value[selected.value].readflag = flag
  }
  const setSeen = (s: { feed_id: number; serial: number }) => {
    const item = list.value.find(
      (e) => e.feed_id === s.feed_id && e.serial === s.serial
    )
    if (item) {
      item.readflag = "Seen"
    }
  }
  const currentEntrySerialData = (): ContentSerialData => {
    const item = list.value[selected.value]
    return {
      serial: item.serial,
      feed_id: item.feed_id,
      readflag: item.readflag,
    }
  }
  const currentEntryURL = (): string | undefined => {
    return list.value[selected.value].url
  }

  const nextContent = () => {
    const index = selected.value + 1
    if (index === list.value.length) {
      return
    }
    selected.value = index
  }
  const previousContent = () => {
    const index = selected.value - 1
    if (0 > index) {
      return
    }
    selected.value = index
  }

  return {
    setContents,
    setIndex,
    setReadflag,
    setSeen,
    list,
    selected,
    currentEntrySerialData,
    currentEntryURL,
    nextContent,
    previousContent,
  }
})
