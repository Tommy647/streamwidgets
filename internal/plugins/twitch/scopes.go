package twitch

const (
	scopeAnalyticsReadExtensions = `analytics:read:extensions`
	scopeAnalyticsReadGames      = `analytics:read:games`

	scopeBitsRead = `bits:read`

	scopeChannelEditCommercial    = `channel:edit:commercial`
	scopeChannelManageBroadcast   = `channel:manage:broadcast`
	scopeChannelManageExtensions  = `channel:manage:extensions`
	scopeChannelManageRedemptions = `channel:manage:redemptions`
	scopeChannelReadHypeTrain     = `channel:read:hype_train`
	scopeChannelReadRedemptions   = `channel:read:redemptions`
	scopeChannelReadStreamKey     = `channel:read:stream_key`
	scopeChannelReadSubscriptions = `channel:read:subscriptions`
	scopeChannelModerate          = `channel:moderate` //  Perform moderation actions in a channel. The user requesting the scope must be a moderator in the channel.

	scopeClipsEdit = `clips:edit`

	scopeUserEdit          = `user:edit`
	scopeUserEditFollows   = `user:edit:follows`
	scopeUserReadBroadcast = `user:read:broadcast`
	scopeUserReadEmail     = `user:read:email`

	scopeChatEdit = `chat:edit` // Send live stream chat and rooms messages.
	scopeChatRead = `chat:read` // View live stream chat and rooms messages.

	scopeWhispersRead = `whispers:read` // View your whisper messages.
	scopeWhispersEdit = `whispers:edit` // Send whisper messages.
)
