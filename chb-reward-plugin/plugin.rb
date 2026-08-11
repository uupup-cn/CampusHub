# name: chb_reward_plugin
# about: 监听 Discourse 事件，转发到 CampusHub 后端奖励引擎
# version: 1.0.0
# authors: CampusHub Team
# url: https://github.com/campushub/chb-reward-plugin

enabled_site_setting :chb_reward_plugin_enabled

register_asset "javascripts/chb_checkin.js"
register_asset "stylesheets/chb_checkin.css"

after_initialize do
  require_dependency "chb_event_handler"
  require_dependency "chb_api_client"
  require_dependency "chb_event_serializer"
  require_dependency "chb_trust_level_sync"

  handler = ChbEventHandler.new

  # 发帖事件
  DiscourseEvent.on(:topic_created) do |topic, opts, user|
    handler.on_topic_created(topic, user)
  end

  # 回复事件
  DiscourseEvent.on(:post_created) do |post, opts, user|
    handler.on_post_created(post, user)
  end

  # 点赞事件
  DiscourseEvent.on(:like_added) do |post, user|
    handler.on_like_added(post, user)
  end

  # 信任等级变更事件
  DiscourseEvent.on(:user_trust_level_change) do |user, old_level, new_level|
    handler.on_trust_level_change(user, new_level)
  end

  # 注册签到路由 - 代理到后端
  add_to_serializer(:current_user, :chb_checkin_url) do
    Discourse.base_url + "/chb/checkin"
  end

  Discourse::Application.routes.append do
    get "/chb/checkin/status" => "chb_checkin#status"
    post "/chb/checkin" => "chb_checkin#checkin"
  end
end
