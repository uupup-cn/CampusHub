class ChbEventHandler
  def initialize
    @client = ChbApiClient::Client.new
  end

  def on_topic_created(topic, user)
    return unless user
    data = ChbEventSerializer.serialize_topic(topic, user)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key], data[:trust_level])
  end

  def on_post_created(post, user)
    return unless user
    # Skip topic first post (topic_created already handles it)
    return if post.post_number == 1
    data = ChbEventSerializer.serialize_post(post, user)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key], data[:trust_level])
  end

  def on_like_added(post, user)
    handle_like(post, user)
  end

  def on_like_created(post, user)
    handle_like(post, user)
  end

  private

  def handle_like(post, user)
    return unless post&.user
    # Don't reward self-likes
    return if post.user_id == user.id
    data = ChbEventSerializer.serialize_like(post, user)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key], data[:trust_level])
  end

  def on_trust_level_change(user, new_level)
    return unless user
    data = ChbEventSerializer.serialize_trust_level(user, new_level)
    @client.sync_trust_level(data[:discourse_user_id], data[:trust_level], data[:idempotency_key])
  end
end
