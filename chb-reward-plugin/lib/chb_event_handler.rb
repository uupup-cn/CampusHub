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

  # like_created event receives a PostAction object
  # post_action.user = the liker
  # post_action.post = the liked post
  # post_action.post.user = the post author (reward recipient)
  def on_like_created(post_action)
    return unless post_action
    post = post_action.post
    liker = post_action.user
    return unless post && liker
    # Don't reward self-likes
    return if post.user_id == liker.id

    data = ChbEventSerializer.serialize_like(post, liker)
    @client.send_reward(data[:action], data[:discourse_user_id], data[:ref_type], data[:ref_id], data[:idempotency_key], data[:trust_level])
  end

  def on_trust_level_change(user, new_level)
    return unless user
    data = ChbEventSerializer.serialize_trust_level(user, new_level)
    @client.sync_trust_level(data[:discourse_user_id], data[:trust_level], data[:idempotency_key])
  end
end
