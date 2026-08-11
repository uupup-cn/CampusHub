u = User.find_by(username: "admin")
if u.nil?
  u = User.new
  u.username = "admin"
  u.email = "admin@campushub.local"
  u.password = "CampusHubAdmin2026!"
  u.admin = true
  u.active = true
  u.approved = true
  u.save!
  puts "Admin created: #{u.username}"
else
  puts "Admin already exists: #{u.username}"
end

SiteSetting.chb_backend_api_url = "http://localhost:9090"
puts "chb_backend_api_url = #{SiteSetting.chb_backend_api_url}"
SiteSetting.chb_backend_api_key = "test-api-key"
puts "chb_backend_api_key = #{SiteSetting.chb_backend_api_key}"
SiteSetting.chb_reward_plugin_enabled = true
puts "chb_reward_plugin_enabled = #{SiteSetting.chb_reward_plugin_enabled}"
SiteSetting.title = "CampusHub Forum"
puts "title = #{SiteSetting.title}"
puts "All settings configured"
