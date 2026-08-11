class ChbTrustLevelSync
  def self.sync_all_users
    client = ChbApiClient::Client.new
    User.human_users.find_each do |user|
      next unless user&id
      data = ChbEventSerializer.serialize_trust_level(user, user.trust_level)
      client.sync_trust_level(data[:discourse_user_id], data[:trust_level], data[:idempotency_key])
    end
  end
end

# 定时同步（每 6 小时）
DiscourseEvent.on(:scheduled_sync) do
  ChbTrustLevelSync.sync_all_users
end
