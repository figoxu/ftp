require 'test_helper'

class HostControllerTest < ActionDispatch::IntegrationTest
  test "should get index" do
    get host_index_url
    assert_response :success
  end

end
