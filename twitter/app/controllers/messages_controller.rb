class MessagesController < ApplicationController
    def initialize
        @result = {
            # params: params,
            "success" => true,
            "error" => {
            }
        }
    end
    
    def index
        conversation = Conversation.find_by(id: params[:id])
        messages = nil
        if (conversation!=nil) then
            messages = conversation.messages
        end
        @result["messages"] = messages

        render json: {result: @result}
    end

    def new
        new_message = Message.new(content: params[:content], user_id: params[:user_id], conversation_id: params[:conversation_id])
        new_message.save
        
        render json: {result: @result}
    end
end
