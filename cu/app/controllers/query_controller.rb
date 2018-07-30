class QueryController < ApplicationController
  
  def index
    result = {
      "params": params,
      "success" => false,
      "error" => {
      },
    }
    render json: {result: result}
  end

  def submit
    result = {
      "params": params,
      "success" => false,
      "error" => {
      },
    }
    result["query_id"] = 18
    render json: {result: result}
  end

  def result
    result = {
      "params": params,
      "success" => false,
      "error" => {
      },
    }
    render json: {result: result}
  end

end
