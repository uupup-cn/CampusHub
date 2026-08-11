class ChbEventHandler
  def initialize
    @client = ChbApiClient::Client.new
  end

  def on_topic_created(topic, user)
    return unless user
    data = ChbEventSerializer.serialize_topic(topic, user)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key])
  end

  def on_post_created(post, user)
    return unless user
    # 跳过主题帖本身（topic_created 已处理）
    return if post.post_number == 1
    data = ChbEventSerializer.serialize_post(post, user)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key])
  end

  def on_like_added(post, user)
    return unless post&.user
    # 不给自己点赞
    return if post.user_id == user.id
    data = ChbEventSerializer.serialize_like(post, user)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key])
  end

  def on_trust_level_change(user, new_level)
    return unless user
    data = ChbEventSerializer.serialize_trust_level(user, new_level)
    @client.sync_trust_level(data[:discourse_user_id], data[:trust_level], data[:idempotency_key])
  end
end
