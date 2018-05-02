class MessagesController < ApplicationController
    def index
        puts "***********"
        puts params[:conversation.id]
        puts "***********"
        render json: {result: @result}
    end
end
