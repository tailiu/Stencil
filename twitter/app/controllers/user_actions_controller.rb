class UserActionsController < ApplicationController
    def index
        puts '*******************'
        puts params[:id]
        puts "gogogogoog"
        puts '*******************'

        if params[:type] == "following_relationship"
            @user = User.find(params[:id])
            @following_num = @user.user_actions.where(action_type: "following").count
            @followed_num = User.joins("INNER JOIN user_actions ON user_actions.to_user_id = users.id").where("user_actions.action_type" => "following").count

            @result = {
                # params: params,
                "followed_num" => @followed_num,
                "following_num" => @following_num
            } 

            render json: {result: @result}
        end

    end
end
