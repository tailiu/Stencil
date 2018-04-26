class UsersController < ApplicationController
    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        @result = {
            # params: params,
        }

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            @result["success"] = true
            @result["user"] = @new_user
        else 
            puts @new_user.errors.messages
            puts @new_credential.errors.messages
            @result["success"] = false
            @result["message"] = "Email already exists!"
        end

        render json: {result: @result}
    end

    def verify
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        if @new_user.valid? && @new_credential.valid?
            @new_user.save
            @new_credential.save
            render json: "something"
            # render "pages/home"
        else 
            puts @new_user.errors.messages
            puts @new_credential.errors.messages
            render json: params
        end
    end

end
