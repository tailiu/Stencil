class UsersController < ApplicationController
    def new
        @new_user = User.new(name: params[:name], handle: params[:handle])
        @new_credential = @new_user.build_credential(email: params[:email], password: params[:password])

        @result = {
            # params: params,
            "success" => false,
            "error" => {
                "message": "",
            }
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
            @result["error"]["message"] = "Email already exists!"
        end

        render json: {result: @result}
    end

    def verify

    end

end
