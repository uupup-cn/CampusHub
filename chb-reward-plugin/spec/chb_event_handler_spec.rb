# frozen_string_literal: true

require 'rails_helper'

describe ChbEventHandler do
  let(:handler) { ChbEventHandler.new }
  let(:user) { Fabricate(:user) }

  context "when plugin is enabled" do
    before do
      SiteSetting.chb_reward_plugin_enabled = true
      SiteSetting.chb_backend_api_url = "http://localhost:8080"
    end

    it "handles topic_created event" do
      topic = Fabricate(:topic, user: user)
      expect { handler.on_topic_created(topic, user) }.not_to raise_error
    end

    it "handles post_created event" do
      post = Fabricate(:post, user: user, post_number: 2)
      expect { handler.on_post_created(post, user) }.not_to raise_error
    end

    it "skips post_number 1 (topic body)" do
      post = Fabricate(:post, user: user, post_number: 1)
      expect(handler.on_post_created(post, user)).to be_nil
    end

    it "handles like_added event" do
      post = Fabricate(:post, user: user)
      liker = Fabricate(:user)
      expect { handler.on_like_added(post, liker) }.not_to raise_error
    end

    it "skips self-like" do
      post = Fabricate(:post, user: user)
      expect(handler.on_like_added(post, user)).to be_nil
    end
  end
end
