class UserActionsController < ApplicationController
    def index
        if params[:type] == "follow"
            @user = User.find(params[:id])
            @following_num = @user.user_actions.where(action_type: "follow").count
            @followed_num = User.joins("INNER JOIN user_actions ON user_actions.to_user_id = users.id").where("user_actions.action_type" => "follow").count

            @result = {
                # params: params,
                "followed_num" => @followed_num,
                "following_num" => @following_num
            }

<<<<<<< HEAD
            @result[:session] = session
=======
>>>>>>> cac4e53ca4506e16e75f317e071093631f3f4b1c
            render json: {result: @result}
        end
    end

    def follow
    end

    def unfollow
    end
    
end
