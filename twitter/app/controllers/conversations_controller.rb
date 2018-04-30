class ConversationsController < ApplicationController
    def new
        
        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
        
        @conversation = Conversation.new()

        render json: {result: @result}
    end
end
