export type Account = {
  id: string;
  username?: string;
  acct?: string;
  display_name?: string;
  url?: string;
  avatar?: string;
  avatar_static?: string;
  header?: string;
  header_static?: string;
  locked?: boolean;
  bot?: boolean;
  created_at?: string;
  note?: string;
};

export type MediaAttachment = {
  id: string;
  type: string;
  url?: string;
  preview_url?: string;
  remote_url?: string | null;
  text_url?: string;
  description?: string | null;
};

export type Mention = {
  id: string;
  username: string;
  url: string;
  acct: string;
};

export type Tag = {
  name: string;
  url: string;
};

export type Application = {
  name?: string;
  website?: string | null;
};
export type Post = {
  id: string;
  created_at: string;
  in_reply_to_id: string | null;
  in_reply_to_account_id: string | null;
  sensitive: boolean;
  spoiler_text: string;
  visibility: 'public' | 'unlisted' | 'private' | 'direct' | string;
  language: string | null;
  uri: string;
  url: string | null;
  replies_count?: number;
  reblogs_count?: number;
  favourites_count?: number;
  favourited?: boolean | null;
  reblogged?: boolean | null;
  muted?: boolean | null;
  bookmarked?: boolean | null;
  pinned?: boolean | null;
  content?: string;
  filtered?: any[];
  account?: Account;
  media_attachments?: MediaAttachment[];
  mentions?: Mention[];
  tags?: Tag[];
  application?: Application | null;
};
