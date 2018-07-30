require 'test_helper'

class QueryControllerTest < ActionDispatch::IntegrationTest
  test "should get index" do
    get query_index_url
    assert_response :success
  end

end
