require 'rails_helper'

describe ChbEventHandler do
  let(:handler) { ChbEventHandler.new }
  let(:client) { double("ChbApiClient::Client") }

  before do
    allow(handler).to receive(:instance_variable_get).with(:@client).and_return(client)
  end

  it "handles like_created event with PostAction" do
    post_action = double("PostAction")
    post = double("Post", id: 1, user_id: 2)
    liker = double("User", id: 3, trust_level: 1)
    allow(post_action).to receive(:post).and_return(post)
    allow(post_action).to receive(:user).and_return(liker)
    allow(client).to receive(:send_reward)

    expect { handler.on_like_created(post_action) }.not_to raise_error
  end

  it "skips self-likes" do
    post_action = double("PostAction")
    post = double("Post", id: 1, user_id: 2)
    liker = double("User", id: 2, trust_level: 1)
    allow(post_action).to receive(:post).and_return(post)
    allow(post_action).to receive(:user).and_return(liker)

    expect(client).not_to receive(:send_reward)
    handler.on_like_created(post_action)
  end
end
