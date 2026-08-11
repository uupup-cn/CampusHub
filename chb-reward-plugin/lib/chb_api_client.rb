module ChbApiClient
  class Client
    def initialize
      @base_url = SiteSetting.chb_backend_api_url
      @api_key = SiteSetting.chb_backend_api_key
    end

    def send_reward(action, user_id, ref_type, ref_id, idempotency_key)
      response = Excon.post(
        "\#{@base_url}/api/chb/reward",
        headers: {
          "Content-Type" => "application/json",
          "X-API-Key" => @api_key
        },
        body: {
          action: action,
          discourse_user_id: user_id,
          ref_type: ref_type,
          ref_id: ref_id,
          idempotency_key: idempotency_key
        }.to_json
      )
      JSON.parse(response.body)
    rescue => e
      Rails.logger.error("CHB Reward API error: \#{e.message}")
      nil
    end

    def sync_trust_level(user_id, trust_level, idempotency_key)
      response = Excon.post(
        "\#{@base_url}/api/chb/sync/trust-level",
        headers: {
          "Content-Type" => "application/json",
          "X-API-Key" => @api_key
        },
        body: {
          discourse_user_id: user_id,
          trust_level: trust_level,
          idempotency_key: idempotency_key
        }.to_json
      )
      JSON.parse(response.body)
    rescue => e
      Rails.logger.error("CHB Trust Level Sync error: \#{e.message}")
      nil
    end

    def checkin(user_id)
      response = Excon.post(
        "\#{@base_url}/api/chb/checkin",
        headers: {
          "Content-Type" => "application/json",
          "X-API-Key" => @api_key
        },
        body: {
          discourse_user_id: user_id
        }.to_json
      )
      JSON.parse(response.body)
    rescue => e
      Rails.logger.error("CHB Checkin API error: \#{e.message}")
      nil
    end

    def checkin_status(user_id)
      response = Excon.get(
        "\#{@base_url}/api/chb/checkin/status?user_id=\#{user_id}",
        headers: {
          "X-API-Key" => @api_key
        }
      )
      JSON.parse(response.body)
    rescue => e
      Rails.logger.error("CHB Checkin Status error: \#{e.message}")
      nil
    end
  end
end
