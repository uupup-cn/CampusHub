class ChbCheckinController < ApplicationController
  before_action :ensure_logged_in

  def checkin
    client = ChbApiClient::Client.new
    result = client.checkin(current_user.id, current_user.trust_level)
    
    if result.nil?
      render json: { code: 500, message: "backend_error" }, status: 500
    else
      render json: result
    end
  end

  def status
    client = ChbApiClient::Client.new
    result = client.checkin_status(current_user.id)
    
    if result.nil?
      render json: { code: 500, message: "backend_error" }, status: 500
    else
      render json: result
    end
  end
end
