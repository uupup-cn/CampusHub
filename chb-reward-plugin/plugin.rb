# name: chb_reward_plugin
# about: Monitor Discourse events and forward to CampusHub backend reward engine
# version: 1.0.0
# authors: CampusHub Team
# url: https://github.com/campushub/chb-reward-plugin

enabled_site_setting :chb_reward_plugin_enabled

register_asset "stylesheets/chb_checkin.css"

after_initialize do
  require_relative "lib/chb_event_handler"
  require_relative "lib/chb_api_client"
  require_relative "lib/chb_event_serializer"
  require_relative "lib/chb_trust_level_sync"
  require_relative "app/controllers/chb_checkin_controller"

  handler = ChbEventHandler.new

  # Topic created event
  DiscourseEvent.on(:topic_created) do |topic, opts, user|
    handler.on_topic_created(topic, user)
  end

  # Post created event (replies)
  DiscourseEvent.on(:post_created) do |post, opts, user|
    handler.on_post_created(post, user)
  end

  # Like events - listen to both like_added and like_created for compatibility
  DiscourseEvent.on(:like_added) do |post, user|
    handler.on_like_added(post, user)
  end

  DiscourseEvent.on(:like_created) do |post, user|
    handler.on_like_created(post, user)
  end

  # Trust level change event
  DiscourseEvent.on(:user_trust_level_change) do |user, old_level, new_level|
    handler.on_trust_level_change(user, new_level)
  end

  # Register checkin URL to serializer
  add_to_serializer(:current_user, :chb_checkin_url) do
    Discourse.base_url + "/chb/checkin"
  end

  Discourse::Application.routes.append do
    get "/chb/checkin/status" => "chb_checkin#status"
    post "/chb/checkin" => "chb_checkin#checkin"
  end
end
