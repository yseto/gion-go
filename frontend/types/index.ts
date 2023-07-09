export type Profile = {
  autoseen: boolean;
  onLoginSkipPinList: boolean;
  entryCount: number;
  substringLength: number;
};

export type PinList = {
  feed_id: number;
  serial: number;
  title: string;
  update_at: string;
  url: string;
};
