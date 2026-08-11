module ChbEventSerializer
  def self.serialize_topic(topic, user)
    {
      action: "post",
      discourse_user_id: user.id,
      ref_type: "topic",
      ref_id: topic.id,
      trust_level: user.trust_level,
      idempotency_key: "topic_#{topic.id}_#{user.id}"
    }
  end

  def self.serialize_post(post, user)
    {
      action: "reply",
      discourse_user_id: user.id,
      ref_type: "post",
      ref_id: post.id,
      trust_level: user.trust_level,
      idempotency_key: "post_#{post.id}_#{user.id}"
    }
  end

  def self.serialize_like(post, user)
    {
      action: "liked",
      discourse_user_id: post.user_id,
      ref_type: "like",
      ref_id: post.id,
      trust_level: user.trust_level,
      idempotency_key: "like_#{post.id}_#{user.id}"
    }
  end

  def self.serialize_trust_level(user, new_level)
    {
      discourse_user_id: user.id,
      trust_level: new_level,
      idempotency_key: "tl_#{user.id}_#{new_level}_#{Time.now.to_i}"
    }
  end
end
